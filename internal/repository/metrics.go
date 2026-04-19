package repository

import (
	"errors"
	"maps"
	"slices"
	"sync"

	"github.com/mikhailpashkov/metrics/internal/model"
)

type MetricsRepository interface {
	FindById(id int64) (*models.Metrics, error)
	FindByName(name string) ([]*models.Metrics, error)
	FindAll() ([]*models.Metrics, error)
	Save(*models.Metrics) (*models.Metrics, error)
	DeleteAll() error
}

type MetricsMemoryRepository struct {
	storage         map[int64]*models.Metrics
	lastId          int64
	mu              sync.RWMutex
	errNotFound     error
	errMetricsIsNil error
}

func NewMetricsMemoryRepository() MetricsRepository {
	return &MetricsMemoryRepository{
		storage:         make(map[int64]*models.Metrics),
		lastId:          -1,
		errNotFound:     errors.New("not found"),
		errMetricsIsNil: errors.New("metrics is nil"),
	}
}

func (m *MetricsMemoryRepository) FindById(id int64) (*models.Metrics, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	metrics, ok := m.storage[id]
	if !ok {
		return nil, m.errNotFound
	}
	return metrics, nil
}

func (m *MetricsMemoryRepository) FindByName(name string) ([]*models.Metrics, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	metrics := make([]*models.Metrics, 0)
	for _, metric := range m.storage {
		if metric.Name == name {
			metrics = append(metrics, metric)
		}
	}
	return metrics, nil
}

func (m *MetricsMemoryRepository) FindAll() ([]*models.Metrics, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return slices.Collect(maps.Values(m.storage)), nil
}

func (m *MetricsMemoryRepository) Save(metrics *models.Metrics) (*models.Metrics, error) {
	if metrics == nil {
		return nil, m.errMetricsIsNil
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if metrics.ID == -1 {
		m.lastId++
		metrics.ID = m.lastId
	}
	m.storage[metrics.ID] = metrics
	return metrics, nil
}

func (m *MetricsMemoryRepository) DeleteAll() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.storage = make(map[int64]*models.Metrics)
	m.lastId = -1
	return nil
}
