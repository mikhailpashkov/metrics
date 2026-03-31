package poller

import models "github.com/mikhailpashkov/metrics/internal/model"

type MetricsPoller interface {
	GetMetrics() ([]*models.Metrics, error)
}
