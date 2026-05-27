package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/mikhailpashkov/metrics/internal/dto"
	"github.com/mikhailpashkov/metrics/internal/mapper"
	models "github.com/mikhailpashkov/metrics/internal/model"
	"github.com/mikhailpashkov/metrics/internal/service"
)

type UpdateMetricsBatchHandler struct {
	logger         *slog.Logger
	metricsService service.MetricsService
}

func NewUpdateMetricsBatchHandlerFunc(logger *slog.Logger, metricsService service.MetricsService) http.HandlerFunc {
	return UpdateMetricsBatchHandler{
		logger:         logger,
		metricsService: metricsService,
	}.serveHttp
}

func (m *UpdateMetricsBatchHandler) serveHttp(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		m.logger.Debug("Method not allowed")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var request dto.UpdateMetricsBatchRequest

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		m.logger.Debug("Error decoding body", "err", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if len(request) == 0 {
		m.logger.Debug("Empty request")
		return
	}

	metricsForUpdate := make([]*models.Metrics, len(request))
	for i, metricsDto := range request {
		if metricsDto.Type == "" {
			m.logger.Debug("Invalid request: Empty type")
			http.Error(w, "Invalid request: Empty type", http.StatusBadRequest)
			return
		}

		if metricsDto.ID == "" {
			m.logger.Debug("Invalid request: Empty name")
			http.Error(w, "Invalid request: Empty name", http.StatusBadRequest)
			return
		}

		metrics := mapper.MetricsFromMetricsDto(metricsDto)
		if !models.IsValidMetrics(metrics) {
			m.logger.Debug("Metric type doesnt match its content", "metrics_id", metricsDto.ID)
			http.Error(
				w, fmt.Sprintf("Metric type doesnt match its content. ID: %s", metricsDto.ID),
				http.StatusBadRequest,
			)
			return
		}
		metricsForUpdate[i] = metrics
	}

	err = m.metricsService.UpdateMetricsBatch(r.Context(), metricsForUpdate)
	if err != nil {
		m.logger.Error("Batch update metrics error", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
