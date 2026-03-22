package service

import (
	"time"

	models "github.com/mikhailpashkov/metrics/internal/model"
	"github.com/mikhailpashkov/metrics/internal/repository"
)

type Metrics struct {
	metricsStorage repository.MetricsStorage
}

func NewMetrics(metricsStorage repository.MetricsStorage) *Metrics {
	return &Metrics{metricsStorage}
}

func (metrics *Metrics) UpdateCounter(name string, delta int64) (*models.Metrics, error) {
	return metrics.updateMetrics(&models.Metrics{
		ID:    -1,
		Type:  models.Counter,
		Name:  name,
		Delta: &delta,
		TS:    time.Now().UnixMilli(),
	})
}

func (metrics *Metrics) UpdateGauge(name string, value float64) (*models.Metrics, error) {
	return metrics.updateMetrics(&models.Metrics{
		ID:    -1,
		Type:  models.Gauge,
		Name:  name,
		Value: &value,
		TS:    time.Now().UnixMilli(),
	})
}

func (metrics *Metrics) updateMetrics(metricsModel *models.Metrics) (*models.Metrics, error) {
	savedMetrics, err := metrics.metricsStorage.Save(metricsModel)
	if err != nil {
		return nil, err
	}

	return savedMetrics, nil
}
