package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"time"

	models "github.com/mikhailpashkov/metrics/internal/model"
	"github.com/mikhailpashkov/metrics/internal/repository"
)

type MetricsService interface {
	UpdateMetrics(ctx context.Context, metricsModel *models.Metrics) (*models.Metrics, error)
	UpdateCounter(ctx context.Context, name string, delta int64) (*models.Metrics, error)
	UpdateGauge(ctx context.Context, name string, value float64) (*models.Metrics, error)
	GetAllRecords(ctx context.Context) ([]*models.Metrics, error)
	GetAllAccumulated(ctx context.Context) ([]*models.Metrics, error)
	DeleteAll(ctx context.Context) error
	Restore(ctx context.Context) error
	SetupBackup(ctx context.Context, storeInterval int) error
}

type MetricsServiceImpl struct {
	logger            *slog.Logger
	metricsRepository repository.MetricsRepository
	backupRepository  repository.BackupRepository
	backupCallback    func(ctx context.Context)
}

func NewMetricsService(logger *slog.Logger, metricsStorage repository.MetricsRepository, backupRepository repository.BackupRepository) MetricsService {
	return &MetricsServiceImpl{
		logger:            logger,
		metricsRepository: metricsStorage,
		backupRepository:  backupRepository,
		backupCallback:    func(ctx context.Context) {},
	}
}

func (ms *MetricsServiceImpl) UpdateMetrics(ctx context.Context, metricsModel *models.Metrics) (*models.Metrics, error) {
	defer ms.backupCallback(ctx)
	savedMetrics, err := ms.metricsRepository.Save(ctx, metricsModel)
	if err != nil {
		return nil, err
	}

	return savedMetrics, nil
}

func (ms *MetricsServiceImpl) UpdateCounter(ctx context.Context, name string, delta int64) (*models.Metrics, error) {
	return ms.UpdateMetrics(ctx, &models.Metrics{
		ID:    -1,
		Type:  models.Counter,
		Name:  name,
		Delta: &delta,
		TS:    time.Now().UnixMilli(),
	})
}

func (ms *MetricsServiceImpl) UpdateGauge(ctx context.Context, name string, value float64) (*models.Metrics, error) {
	return ms.UpdateMetrics(ctx, &models.Metrics{
		ID:    -1,
		Type:  models.Gauge,
		Name:  name,
		Value: &value,
		TS:    time.Now().UnixMilli(),
	})
}

func (ms *MetricsServiceImpl) GetAllRecords(ctx context.Context) ([]*models.Metrics, error) {
	return ms.metricsRepository.FindAll(ctx)
}

func (ms *MetricsServiceImpl) GetAllAccumulated(ctx context.Context) ([]*models.Metrics, error) {
	records, err := ms.GetAllRecords(ctx)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return []*models.Metrics{}, nil
	}

	nameToRecords := make(map[string][]*models.Metrics)
	for _, record := range records {
		_, ok := nameToRecords[record.Name]
		if !ok {
			nameToRecords[record.Name] = make([]*models.Metrics, 0)
		}
		nameToRecords[record.Name] = append(nameToRecords[record.Name], record)
	}

	result := make([]*models.Metrics, 0)
	for name, groupedRecords := range nameToRecords {
		recordsType := groupedRecords[0].Type
		for _, record := range groupedRecords {
			if record.Type != recordsType {
				return nil, fmt.Errorf("record type mismatch")
			}
		}

		switch recordsType {
		case models.Counter:
			var accumulatedDelta int64
			for _, record := range groupedRecords {
				if record.Delta == nil {
					ms.logger.Warn("counter delta is nil", "id", record.ID, "name", record.Name)
					continue
				}
				accumulatedDelta += *record.Delta
			}
			accumulatedMetric := &models.Metrics{
				ID:    -1,
				Type:  models.Counter,
				Name:  name,
				Delta: &accumulatedDelta,
				Value: nil,
				TS:    0,
			}
			result = append(result, accumulatedMetric)
			continue
		case models.Gauge:
			sort.Slice(groupedRecords, func(i, j int) bool {
				return groupedRecords[i].ID < groupedRecords[j].ID
			})
			lastRecordByID := groupedRecords[len(groupedRecords)-1]

			result = append(result, lastRecordByID)
			continue
		default:
			return nil, fmt.Errorf("invalid record type")
		}
	}

	return result, nil
}

func (ms *MetricsServiceImpl) DeleteAll(ctx context.Context) error {
	defer ms.backupCallback(ctx)
	return ms.metricsRepository.DeleteAll(ctx)
}

func (ms *MetricsServiceImpl) Restore(ctx context.Context) error {
	restoredMetrics, err := ms.backupRepository.FindAll(ctx)
	if err != nil {
		return err
	}

	savedIds := make([]int64, 0, len(restoredMetrics))
	errs := make([]error, 0)
	for _, bMetrics := range restoredMetrics {
		toSave := &models.Metrics{
			ID:    -1,
			Type:  bMetrics.Type,
			Name:  bMetrics.ID,
			Delta: bMetrics.Delta,
			Value: bMetrics.Value,
			TS:    time.Now().UnixMilli(),
		}
		saved, ferr := ms.metricsRepository.Save(ctx, toSave)
		if ferr != nil {
			errs = append(errs, ferr)
			continue
		}
		savedIds = append(savedIds, saved.ID)
		ms.logger.Debug("Restored",
			"id", saved.ID,
			"type", saved.Type,
			"name", saved.Name,
		)
	}

	if len(errs) > 0 {
		for _, id := range savedIds {
			ferr := ms.metricsRepository.DeleteById(ctx, id)
			if ferr != nil {
				errs = append(errs, ferr)
			}
		}
		return errors.Join(errs...)
	}

	return nil
}

func (ms *MetricsServiceImpl) SetupBackup(ctx context.Context, storeInterval int) error {
	if storeInterval == 0 {
		ms.backupCallback = ms.doBackup
		return nil
	}

	go func() {
		ticker := time.NewTicker(time.Duration(storeInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				ms.logger.InfoContext(ctx, "backup stopped", "why", ctx.Err())
				return
			case <-ticker.C:
				ms.doBackup(ctx)
			}
		}
	}()

	return nil
}

func (ms *MetricsServiceImpl) doBackup(ctx context.Context) {
	ms.logger.Info("backup started")
	metrics, err := ms.GetAllAccumulated(ctx)
	if err != nil {
		ms.logger.Error("get metrics failed", "err", err)
		return
	}

	backupMetrics := make([]*models.BackupMetrics, 0, len(metrics))
	for _, m := range metrics {
		backupMetrics = append(backupMetrics, &models.BackupMetrics{
			ID:    m.Name,
			Type:  m.Type,
			Delta: m.Delta,
			Value: m.Value,
		})
	}

	err = ms.backupRepository.SaveAll(ctx, backupMetrics)
	if err != nil {
		ms.logger.Error("save metrics failed", "err", err)
		return
	}
	ms.logger.Debug("backup finished")
}
