package repository

import (
	"context"
	"maps"
	"slices"
	"sync"

	"github.com/mikhailpashkov/metrics/internal/model"
)

type MetricsMemoryRepository struct {
	storage map[int64]*models.Metrics
	lastId  int64
	mu      sync.RWMutex
}

func NewMetricsMemoryRepository() *MetricsMemoryRepository {
	return &MetricsMemoryRepository{
		storage: make(map[int64]*models.Metrics),
		lastId:  -1,
	}
}

func (r *MetricsMemoryRepository) FindById(ctx context.Context, id int64) (*models.Metrics, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	metrics, ok := r.storage[id]
	if !ok {
		return nil, models.ErrNotFound
	}
	return metrics, nil
}

func (r *MetricsMemoryRepository) FindByName(ctx context.Context, name string) ([]*models.Metrics, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	metrics := make([]*models.Metrics, 0)
	for _, metric := range r.storage {
		if metric.Name == name {
			metrics = append(metrics, metric)
		}
	}
	return metrics, nil
}

func (r *MetricsMemoryRepository) FindAll(ctx context.Context) ([]*models.Metrics, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return slices.Collect(maps.Values(r.storage)), nil
}

func (r *MetricsMemoryRepository) Save(ctx context.Context, metrics *models.Metrics) (*models.Metrics, error) {
	if metrics == nil {
		return nil, models.ErrMetricsIsNil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if metrics.ID == models.MetricsNewID {
		r.lastId++
		metrics.ID = r.lastId
	}
	r.storage[metrics.ID] = metrics
	return metrics, nil
}

func (r *MetricsMemoryRepository) InsertBatch(ctx context.Context, metrics []*models.Metrics) error {
	if metrics == nil {
		return models.ErrMetricsIsNil
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	for _, m := range metrics {
		r.lastId++
		m.ID = r.lastId
		r.storage[m.ID] = m
	}
	return nil
}

func (r *MetricsMemoryRepository) DeleteAll(ctx context.Context) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.storage = make(map[int64]*models.Metrics)
	r.lastId = -1
	return nil
}

func (r *MetricsMemoryRepository) DeleteById(ctx context.Context, id int64) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.storage, id)
	return nil
}
