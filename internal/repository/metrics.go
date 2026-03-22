package repository

import (
	"errors"
	"sync"

	"github.com/mikhailpashkov/metrics/internal/model"
)

type MetricsStorage interface {
	FindById(id int64) (*models.Metrics, error)
	FindByName(name string) ([]models.Metrics, error)
	Save(*models.Metrics) (*models.Metrics, error)
}

type MetricsMemoryStorage struct {
	storage map[int64]*models.Metrics
	lastId  int64
	op      sync.Mutex
}

func NewMetricsMemoryStorage() MetricsStorage {
	return &MetricsMemoryStorage{
		storage: make(map[int64]*models.Metrics),
		lastId:  -1,
	}
}

func (m *MetricsMemoryStorage) FindById(id int64) (*models.Metrics, error) {
	m.op.Lock()
	defer m.op.Unlock()

	metrics, ok := m.storage[id]
	if !ok {
		return nil, errors.New("not found")
	}
	return metrics, nil
}

func (m *MetricsMemoryStorage) FindByName(name string) ([]models.Metrics, error) {
	m.op.Lock()
	defer m.op.Unlock()

	metrics := make([]models.Metrics, 0)
	for _, metric := range m.storage {
		if metric.Name == name {
			metrics = append(metrics, *metric)
		}
	}
	return metrics, nil
}

func (m *MetricsMemoryStorage) Save(metrics *models.Metrics) (*models.Metrics, error) {
	if metrics == nil {
		return nil, errors.New("metrics is nil")
	}

	m.op.Lock()
	defer m.op.Unlock()

	if metrics.ID == -1 {
		m.lastId++
		metrics.ID = m.lastId
	}
	m.storage[metrics.ID] = metrics
	return metrics, nil
}
