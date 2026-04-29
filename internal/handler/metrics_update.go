package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/mikhailpashkov/metrics/internal/dto"
	"github.com/mikhailpashkov/metrics/internal/mapper"
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

func (m *UpdateMetricsHandler) GetUrlPattern() string {
	return "/update"
}

func (m *UpdateMetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	metrics, err := mapper.MetricsFromUpdateMetricsRequest(request)
	if errors.Is(err, strconv.ErrSyntax) {
		http.Error(w, "Invalid request: Invalid value", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "Can't map MetricsFromUpdateMetricsRequest: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = m.metricsService.UpdateMetrics(r.Context(), metrics)
	if err != nil {
		http.Error(w, fmt.Sprintf("Value update error: %s", err), http.StatusInternalServerError)
		return
	}
}
