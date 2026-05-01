package reporter

import (
	"encoding/json"
	"fmt"
	"log/slog"

	models "github.com/mikhailpashkov/metrics/internal/model"
)

type LogReporter struct {
	logger *slog.Logger
}

func NewLogReporter(logger *slog.Logger) *LogReporter {
	return &LogReporter{logger: logger}
}

func (r *LogReporter) GetLogger() *slog.Logger {
	return r.logger
}

func (r *LogReporter) SendMetrics(metrics *models.Metrics) error {
	marshal, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("failed to marshal metrics: %w", err)
	}
	r.GetLogger().Info("Send metrics to log..", "payload", string(marshal))
	return nil
}
