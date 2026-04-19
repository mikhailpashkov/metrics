package service

import (
	"context"
	"fmt"
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
}

type MetricsServiceImpl struct {
	metricsRepository repository.MetricsRepository
}

func NewMetricsService(metricsStorage repository.MetricsRepository) *MetricsServiceImpl {
	return &MetricsServiceImpl{metricsStorage}
}

func (ms *MetricsServiceImpl) UpdateMetrics(ctx context.Context, metricsModel *models.Metrics) (*models.Metrics, error) {
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
				panic("record type mismatch")
			}
		}

		switch recordsType {
		case models.Counter:
			var accumulatedDelta int64
			for _, record := range groupedRecords {
				if record.Delta == nil {
					fmt.Println("[ERR] counter delta is nil")
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
			lastRecordByTS := groupedRecords[len(groupedRecords)-1]

			result = append(result, lastRecordByTS)
			continue
		default:
			panic("invalid record type")
		}
	}

	return result, nil
}

func (ms *MetricsServiceImpl) DeleteAll(ctx context.Context) error {
	return ms.metricsRepository.DeleteAll(ctx)
}
