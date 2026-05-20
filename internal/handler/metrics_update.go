package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/mikhailpashkov/metrics/internal/dto"
	"github.com/mikhailpashkov/metrics/internal/mapper"
	models "github.com/mikhailpashkov/metrics/internal/model"
	"github.com/mikhailpashkov/metrics/internal/service"
)

type UpdateMetricsHandler struct {
	logger         *slog.Logger
	metricsService service.MetricsService
}

func NewUpdateMetricsHandler(logger *slog.Logger, metricsService service.MetricsService) *UpdateMetricsHandler {
	return &UpdateMetricsHandler{
		logger:         logger,
		metricsService: metricsService,
	}
}

func (m *UpdateMetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		m.logger.Debug("Method not allowed")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var request dto.UpdateMetricsRequest

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		m.logger.Debug("Error decoding body", "err", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if request.Type == "" {
		m.logger.Debug("Invalid request: Empty type")
		http.Error(w, "Invalid request: Empty type", http.StatusBadRequest)
		return
	}

	if request.ID == "" {
		m.logger.Debug("Invalid request: Empty name")
		http.Error(w, "Invalid request: Empty name", http.StatusBadRequest)
		return
	}

	metrics := mapper.MetricsFromUpdateMetricsRequest(request)

	isValid := models.IsValidMetrics(metrics)
	if !isValid {
		m.logger.Debug("Metric type doesnt match its content")
		http.Error(w, "Metric type doesnt match its content", http.StatusBadRequest)
		return
	}

	_, err = m.metricsService.UpdateMetrics(r.Context(), metrics)
	if err != nil {
		m.logger.Error("Value update error", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	_, err = w.Write([]byte("{}"))
	if err != nil {
		m.logger.Error("failed to write response", "err", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
