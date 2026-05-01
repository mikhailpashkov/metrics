package service

import (
	"context"
	"fmt"
	"log/slog"
	"sort"
	"time"

	"github.com/google/uuid"
	models "github.com/mikhailpashkov/metrics/internal/model"
	"github.com/mikhailpashkov/metrics/internal/repository"
)

type MetricsService interface {
	UpdateMetrics(ctx context.Context, metricsModel *models.Metrics) (*models.Metrics, error)
	UpdateCounter(ctx context.Context, name string, delta int64) (*models.Metrics, error)
	UpdateGauge(ctx context.Context, name string, value float64) (*models.Metrics, error)
	GetAllRecords(ctx context.Context) ([]*models.Metrics, error)
	GetAllAccumulated(ctx context.Context) ([]*models.Metrics, error)
	Delete(ctx context.Context, id int64) error
	DeleteAll(ctx context.Context) error
}

type MetricsServiceImpl struct {
	logger            *slog.Logger
	metricsRepository repository.MetricsRepository
	eventService      EventService
}

func NewMetricsService(logger *slog.Logger, metricsRepository repository.MetricsRepository, eventService EventService) MetricsService {
	return &MetricsServiceImpl{
		logger:            logger,
		metricsRepository: metricsRepository,
		eventService:      eventService,
	}
}

func (ms *MetricsServiceImpl) UpdateMetrics(ctx context.Context, metricsModel *models.Metrics) (*models.Metrics, error) {
	savedMetrics, err := ms.metricsRepository.Save(ctx, metricsModel)
	if err != nil {
		return nil, err
	}

	defer ms.eventService.Notify(&models.Event{ID: uuid.NewString(), Key: models.MetricsUpdatedEvent})

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

func (ms *MetricsServiceImpl) Delete(ctx context.Context, id int64) error {
	err := ms.metricsRepository.DeleteById(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete metrics record with id %d: %w", id, err)
	}

	defer ms.eventService.Notify(&models.Event{ID: uuid.NewString(), Key: models.MetricsDeletedEvent})

	return nil
}

func (ms *MetricsServiceImpl) DeleteAll(ctx context.Context) error {
	err := ms.metricsRepository.DeleteAll(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete all metrics: %w", err)
	}

	defer ms.eventService.Notify(&models.Event{ID: uuid.NewString(), Key: models.MetricsDeletedEvent})

	return nil
}
