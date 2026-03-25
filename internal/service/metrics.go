package service

import (
	"time"

	models "github.com/mikhailpashkov/metrics/internal/model"
	"github.com/mikhailpashkov/metrics/internal/repository"
)

type MetricsService struct {
	metricsRepository repository.MetricsRepository
}

func NewMetricsService(metricsStorage repository.MetricsRepository) *MetricsService {
	return &MetricsService{metricsStorage}
}

func (ms *MetricsService) UpdateCounter(name string, delta int64) (*models.Metrics, error) {
	return ms.updateMetrics(&models.Metrics{
		ID:    -1,
		Type:  models.Counter,
		Name:  name,
		Delta: &delta,
		TS:    time.Now().UnixMilli(),
	})
}

func (ms *MetricsService) UpdateGauge(name string, value float64) (*models.Metrics, error) {
	return ms.updateMetrics(&models.Metrics{
		ID:    -1,
		Type:  models.Gauge,
		Name:  name,
		Value: &value,
		TS:    time.Now().UnixMilli(),
	})
}

func (ms *MetricsService) GetAllRecords() ([]*models.Metrics, error) {
	return ms.metricsRepository.FindAll()
}

func (ms *MetricsService) updateMetrics(metricsModel *models.Metrics) (*models.Metrics, error) {
	savedMetrics, err := ms.metricsRepository.Save(metricsModel)
	if err != nil {
		return nil, err
	}

	return savedMetrics, nil
}
