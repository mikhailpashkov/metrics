package repository

import (
	"context"
	"fmt"
	"log/slog"

	metricsdb "github.com/mikhailpashkov/metrics/db/metrics"
	"github.com/mikhailpashkov/metrics/internal/mapper"
	models "github.com/mikhailpashkov/metrics/internal/model"
)

type MetricsDBRepository struct {
	metricsQuery *metricsdb.Queries
	logger       *slog.Logger
}

func NewMetricsDBRepository(metricsQuery *metricsdb.Queries, logger *slog.Logger) *MetricsDBRepository {
	return &MetricsDBRepository{
		metricsQuery: metricsQuery,
		logger:       logger,
	}
}

func (r *MetricsDBRepository) FindById(ctx context.Context, id int64) (*models.Metrics, error) {
	metrics, err := r.metricsQuery.FindById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to FindById metrics: %w", err)
	}
	return mapper.MetricsFromDB(&metrics), nil
}

func (r *MetricsDBRepository) FindByName(ctx context.Context, name string) ([]*models.Metrics, error) {
	metrics, err := r.metricsQuery.FindByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to FindByName metrics: %w", err)
	}
	return mapper.MetricsFromDBList(metrics), nil
}

func (r *MetricsDBRepository) FindAll(ctx context.Context) ([]*models.Metrics, error) {
	metrics, err := r.metricsQuery.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to FindAll metrics: %w", err)
	}
	return mapper.MetricsFromDBList(metrics), nil
}

func (r *MetricsDBRepository) Save(ctx context.Context, metrics *models.Metrics) (*models.Metrics, error) {
	if metrics == nil {
		return nil, models.ErrMetricsIsNil
	}

	var result metricsdb.Metric
	err := PGRetry(ctx, func() error {
		var e error
		isExist := metrics.ID != models.MetricsNewID
		if isExist {
			result, e = r.metricsQuery.Update(ctx, *mapper.MetricsToDBUpdateParams(metrics))
		} else {
			result, e = r.metricsQuery.Insert(ctx, *mapper.MetricsToDBInsertParams(metrics))
		}
		return e
	}, r.logger)

	if err != nil {
		return nil, fmt.Errorf("failed to Save metrics: %w", err)
	}
	return mapper.MetricsFromDB(&result), nil
}

func (r *MetricsDBRepository) InsertBatch(ctx context.Context, metrics []*models.Metrics) error {
	if metrics == nil {
		return models.ErrMetricsIsNil
	}
	if len(metrics) == 0 {
		return nil
	}
	for _, m := range metrics {
		if m == nil {
			return models.ErrMetricsIsNil
		}
	}

	rows := mapper.MetricsToDBInsertBatchParamsList(metrics)
	err := PGRetry(ctx, func() error {
		_, err := r.metricsQuery.InsertBatch(ctx, rows)
		return err
	}, r.logger)

	if err != nil {
		return fmt.Errorf("failed to InsertBatch metrics: %w", err)
	}
	return nil
}

func (r *MetricsDBRepository) DeleteAll(ctx context.Context) error {
	return models.ErrDeleteAllNotAllowed
}

func (r *MetricsDBRepository) DeleteById(ctx context.Context, id int64) error {
	var rowsAffected int64
	err := PGRetry(ctx, func() error {
		var e error
		rowsAffected, e = r.metricsQuery.DeleteById(ctx, id)
		return e
	}, r.logger)

	if err != nil {
		return fmt.Errorf("failed to DeleteById metrics: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("not found metrics. id %d: %w", id, models.ErrNotFound)
	}
	return nil
}
