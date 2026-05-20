package repository

import (
	"context"
	"errors"
	"fmt"

	metricsdb "github.com/mikhailpashkov/metrics/db/metrics"
	"github.com/mikhailpashkov/metrics/internal/mapper"
	models "github.com/mikhailpashkov/metrics/internal/model"
)

type MetricsDBRepository struct {
	metricsQuery    *metricsdb.Queries
	errMetricsIsNil error
}

func NewMetricsDBRepository(metricsQuery *metricsdb.Queries) *MetricsDBRepository {
	return &MetricsDBRepository{
		metricsQuery:    metricsQuery,
		errMetricsIsNil: errors.New("metrics is nil"),
	}
}

func (r MetricsDBRepository) FindById(ctx context.Context, id int64) (*models.Metrics, error) {
	metrics, err := r.metricsQuery.FindById(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to FindById metrics: %w", err)
	}
	return mapper.MetricsFromDB(&metrics), nil
}

func (r MetricsDBRepository) FindByName(ctx context.Context, name string) ([]*models.Metrics, error) {
	metrics, err := r.metricsQuery.FindByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("failed to FindByName metrics: %w", err)
	}
	return mapper.MetricsFromDBList(metrics), nil
}

func (r MetricsDBRepository) FindAll(ctx context.Context) ([]*models.Metrics, error) {
	metrics, err := r.metricsQuery.FindAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to FindAll metrics: %w", err)
	}
	return mapper.MetricsFromDBList(metrics), nil
}

func (r MetricsDBRepository) Save(ctx context.Context, metrics *models.Metrics) (*models.Metrics, error) {
	if metrics == nil {
		return nil, r.errMetricsIsNil
	}

	var result metricsdb.Metric
	var err error
	isExist := metrics.ID != -1
	if isExist {
		result, err = r.metricsQuery.Update(ctx, *mapper.MetricsToDBUpdateParams(metrics))
	} else {
		result, err = r.metricsQuery.Insert(ctx, *mapper.MetricsToDBInsertParams(metrics))
	}
	if err != nil {
		return nil, fmt.Errorf("failed to Save metrics: %w", err)
	}
	return mapper.MetricsFromDB(&result), nil
}

func (r MetricsDBRepository) InsertBatch(ctx context.Context, metrics []*models.Metrics) error {
	if metrics == nil {
		return r.errMetricsIsNil
	}
	for _, m := range metrics {
		if m == nil {
			return r.errMetricsIsNil
		}
	}

	rows := mapper.MetricsToDBInsertBatchParamsList(metrics)
	_, err := r.metricsQuery.InsertBatch(ctx, rows)
	if err != nil {
		return fmt.Errorf("failed to InsertBatch metrics: %w", err)
	}
	return nil
}

func (r MetricsDBRepository) DeleteAll(ctx context.Context) error {
	return errors.New("DeleteAll not allowed")
}

func (r MetricsDBRepository) DeleteById(ctx context.Context, id int64) error {
	rowsAffected, err := r.metricsQuery.DeleteById(ctx, id)
	if rowsAffected == 0 {
		return fmt.Errorf("no metrics with id %d was found", id)
	}
	if err != nil {
		return fmt.Errorf("failed to DeleteById metrics: %w", err)
	}
	return nil
}
