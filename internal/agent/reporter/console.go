package reporter

import (
	"encoding/json"
	"log/slog"

	models "github.com/mikhailpashkov/metrics/internal/model"
)

type ConsoleReporter struct {
	logger *slog.Logger
}

func NewConsoleReporter(logger *slog.Logger) MetricsReporter {
	return &ConsoleReporter{logger: logger}
}

func (r *ConsoleReporter) GetLogger() *slog.Logger {
	return r.logger
}

func (r *ConsoleReporter) SendMetrics(metrics *models.Metrics) error {
	marshal, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	r.GetLogger().Info(string(marshal))
	return nil
}
