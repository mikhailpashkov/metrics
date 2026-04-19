package reporter

import models "github.com/mikhailpashkov/metrics/internal/model"

type MetricsReporter interface {
	SendMetrics(metrics *models.Metrics) error
}
