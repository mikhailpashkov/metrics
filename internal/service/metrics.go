package service

import (
	"fmt"
	"sort"
	"time"

	models "github.com/mikhailpashkov/metrics/internal/model"
	"github.com/mikhailpashkov/metrics/internal/repository"
)

type MetricsService interface {
	UpdateMetrics(metricsModel *models.Metrics) (*models.Metrics, error)
	UpdateCounter(name string, delta int64) (*models.Metrics, error)
	UpdateGauge(name string, value float64) (*models.Metrics, error)
	GetAllRecords() ([]*models.Metrics, error)
	GetAllAccumulated() ([]*models.Metrics, error)
	DeleteAll() error
}

type MetricsServiceImpl struct {
	metricsRepository repository.MetricsRepository
}

func NewMetricsService(metricsStorage repository.MetricsRepository) *MetricsServiceImpl {
	return &MetricsServiceImpl{metricsStorage}
}

func (ms *MetricsServiceImpl) UpdateMetrics(metricsModel *models.Metrics) (*models.Metrics, error) {
	savedMetrics, err := ms.metricsRepository.Save(metricsModel)
	if err != nil {
		return nil, err
	}

	return savedMetrics, nil
}

func (ms *MetricsServiceImpl) UpdateCounter(name string, delta int64) (*models.Metrics, error) {
	return ms.UpdateMetrics(&models.Metrics{
		ID:    -1,
		Type:  models.Counter,
		Name:  name,
		Delta: &delta,
		TS:    time.Now().UnixMilli(),
	})
}

func (ms *MetricsServiceImpl) UpdateGauge(name string, value float64) (*models.Metrics, error) {
	return ms.UpdateMetrics(&models.Metrics{
		ID:    -1,
		Type:  models.Gauge,
		Name:  name,
		Value: &value,
		TS:    time.Now().UnixMilli(),
	})
}

func (ms *MetricsServiceImpl) GetAllRecords() ([]*models.Metrics, error) {
	return ms.metricsRepository.FindAll()
}

func (ms *MetricsServiceImpl) GetAllAccumulated() ([]*models.Metrics, error) {
	records, err := ms.GetAllRecords()
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

		if recordsType == models.Counter {
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
		}

		if recordsType == models.Gauge {
			sort.Slice(groupedRecords, func(i, j int) bool {
				return groupedRecords[i].TS > groupedRecords[j].TS
			})
			lastRecordByTS := groupedRecords[len(groupedRecords)-1]

			result = append(result, lastRecordByTS)
			continue
		}

		panic("invalid record type")
	}

	return result, nil
}

func (ms *MetricsServiceImpl) DeleteAll() error {
	return ms.metricsRepository.DeleteAll()
}
