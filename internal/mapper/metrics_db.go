package mapper

import (
	"time"

	models "github.com/mikhailpashkov/metrics/internal/model"
)
import metricsdb "github.com/mikhailpashkov/metrics/db/metrics"

func MetricsToDBInsertParams(metrics *models.Metrics) *metricsdb.InsertParams {
	return &metricsdb.InsertParams{
		Ts:    time.UnixMilli(metrics.TS),
		Type:  metrics.Type,
		Name:  metrics.Name,
		Delta: metrics.Delta,
		Value: metrics.Value,
	}
}

func MetricsToDBUpdateParams(metrics *models.Metrics) *metricsdb.UpdateParams {
	return &metricsdb.UpdateParams{
		ID:    metrics.ID,
		Ts:    time.UnixMilli(metrics.TS),
		Type:  metrics.Type,
		Name:  metrics.Name,
		Delta: metrics.Delta,
		Value: metrics.Value,
	}
}

func MetricsFromDB(metrics *metricsdb.Metric) *models.Metrics {
	return &models.Metrics{
		ID:    metrics.ID,
		Type:  metrics.Type,
		Name:  metrics.Name,
		Delta: metrics.Delta,
		Value: metrics.Value,
		TS:    metrics.Ts.UnixMilli(),
	}
}

func MetricsFromDBList(metrics []metricsdb.Metric) []*models.Metrics {
	result := make([]*models.Metrics, len(metrics))
	for i, m := range metrics {
		result[i] = MetricsFromDB(&m)
	}
	return result
}
