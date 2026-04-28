package reporter

import (
	"log/slog"

	models "github.com/mikhailpashkov/metrics/internal/model"
)

type MetricsReporter interface {
	SendMetrics(metrics *models.Metrics) error
	GetLogger() *slog.Logger
}
