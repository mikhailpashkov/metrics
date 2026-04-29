package repository

import (
	"context"
	"errors"
	"maps"
	"slices"
	"sync"

	"github.com/mikhailpashkov/metrics/internal/model"
)

type MetricsRepository interface {
	FindById(ctx context.Context, id int64) (*models.Metrics, error)
	FindByName(ctx context.Context, name string) ([]*models.Metrics, error)
	FindAll(ctx context.Context) ([]*models.Metrics, error)
	Save(ctx context.Context, metrics *models.Metrics) (*models.Metrics, error)
	DeleteAll(ctx context.Context) error
	DeleteById(ctx context.Context, id int64) error
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

func (m *MetricsMemoryRepository) FindById(ctx context.Context, id int64) (*models.Metrics, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	metrics, ok := m.storage[id]
	if !ok {
		return nil, m.errNotFound
	}
	return metrics, nil
}

func (m *MetricsMemoryRepository) FindByName(ctx context.Context, name string) ([]*models.Metrics, error) {
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

func (m *MetricsMemoryRepository) FindAll(ctx context.Context) ([]*models.Metrics, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return slices.Collect(maps.Values(m.storage)), nil
}

func (m *MetricsMemoryRepository) Save(ctx context.Context, metrics *models.Metrics) (*models.Metrics, error) {
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

func (m *MetricsMemoryRepository) DeleteAll(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.storage = make(map[int64]*models.Metrics)
	m.lastId = -1
	return nil
}

func (m *MetricsMemoryRepository) DeleteById(ctx context.Context, id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	delete(m.storage, id)
	return nil
}
