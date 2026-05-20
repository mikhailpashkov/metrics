package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	models "github.com/mikhailpashkov/metrics/internal/model"
)

type BackupService interface {
	Restore(ctx context.Context) error
	SetupBackup(ctx context.Context, storeInterval int) error
}

type BackupRepository interface {
	FindAll(ctx context.Context) ([]*models.BackupMetrics, error)
	SaveAll(ctx context.Context, metrics []*models.BackupMetrics) error
}

type BackupServiceImpl struct {
	logger           *slog.Logger
	metricsService   MetricsService
	backupRepository BackupRepository
	eventService     EventService
}

func NewBackupService(
	logger *slog.Logger,
	metricsService MetricsService,
	backupRepository BackupRepository,
	eventService EventService,
) *BackupServiceImpl {
	return &BackupServiceImpl{
		logger:           logger,
		metricsService:   metricsService,
		backupRepository: backupRepository,
		eventService:     eventService,
	}
}

func (bs *BackupServiceImpl) Restore(ctx context.Context) error {
	restoredMetrics, err := bs.backupRepository.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to get backup records: %w", err)
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
		saved, ferr := bs.metricsService.UpdateMetrics(ctx, toSave)
		if ferr != nil {
			errs = append(errs, ferr)
			continue
		}
		savedIds = append(savedIds, saved.ID)
		bs.logger.Debug("Restored",
			"id", saved.ID,
			"type", saved.Type,
			"name", saved.Name,
		)
	}

	if len(errs) > 0 {
		for _, id := range savedIds {
			ferr := bs.metricsService.Delete(ctx, id)
			if ferr != nil {
				errs = append(errs, ferr)
			}
		}
		return errors.Join(errs...)
	}

	return nil
}

func (bs *BackupServiceImpl) SetupBackup(ctx context.Context, storeInterval int) error {
	if storeInterval < 0 {
		return errors.New("store interval cannot be negative")
	}
	if storeInterval == 0 {
		bs.eventDrivenBackup(ctx)
		return nil
	}
	bs.periodicBackup(ctx, storeInterval)
	return nil
}

func (bs *BackupServiceImpl) eventDrivenBackup(ctx context.Context) {
	subscriber := &models.EventSubscriber{
		Name:     "backupService",
		Callback: func(_ models.Event) { bs.doBackup(ctx) },
	}

	bs.eventService.Subscribe(models.MetricsDeletedEvent, subscriber)
	bs.eventService.Subscribe(models.MetricsUpdatedEvent, subscriber)
}

func (bs *BackupServiceImpl) periodicBackup(ctx context.Context, storeInterval int) {
	go func() {
		ticker := time.NewTicker(time.Duration(storeInterval) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				bs.logger.InfoContext(ctx, "backup stopped", "why", ctx.Err())
				return
			case <-ticker.C:
				bs.doBackup(ctx)
			}
		}
	}()
}

func (bs *BackupServiceImpl) doBackup(ctx context.Context) {
	bs.logger.Info("backup started")
	metrics, err := bs.metricsService.GetAllAccumulated(ctx)
	if err != nil {
		bs.logger.Error("Error getting all accumulated metrics", "err", err)
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

	err = bs.backupRepository.SaveAll(ctx, backupMetrics)
	if err != nil {
		bs.logger.Error("save metrics failed", "err", err)
		return
	}
	bs.logger.Debug("backup finished")
}
