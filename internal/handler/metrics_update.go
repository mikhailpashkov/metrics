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

func (m *UpdateMetricsHandler) GetLogger() *slog.Logger { return m.logger }

func (m *UpdateMetricsHandler) GetUrlPatterns() []string {
	return []string{"/update", "/update/"}
}

func (m *UpdateMetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var request dto.UpdateMetricsRequest

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if request.Type == "" {
		http.Error(w, "Invalid request: Empty type", http.StatusBadRequest)
		return
	}

	if request.ID == "" {
		http.Error(w, "Invalid request: Empty name", http.StatusBadRequest)
		return
	}

	metrics := mapper.MetricsFromUpdateMetricsRequest(request)

	isValid := models.IsValidMetrics(metrics)
	if !isValid {
		http.Error(w, "Metric type doesnt match its content", http.StatusBadRequest)
		return
	}

	_, err = m.metricsService.UpdateMetrics(r.Context(), metrics)
	if err != nil {
		http.Error(w, fmt.Sprintf("Value update error: %s", err), http.StatusInternalServerError)
		return
	}
}
