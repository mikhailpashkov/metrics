package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

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

func (m *UpdateMetricsHandler) GetUrlPattern() string {
	return "/update/{type}/{name}/{value}"
}

func (m *UpdateMetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	mType := r.PathValue("type")
	name := r.PathValue("name")
	valueStr := r.PathValue("value")

	if mType == "" {
		http.Error(w, "Empty type", http.StatusBadRequest)
		return
	}

	if name == "" {
		http.Error(w, "Empty name", http.StatusBadRequest)
		return
	}

	switch mType {
	case models.Counter:
		parseInt, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("Value conversion error: %s", err), http.StatusBadRequest)
			return
		}

		_, err = m.metricsService.UpdateCounter(r.Context(), name, parseInt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Value update error: %s", err), http.StatusInternalServerError)
			return
		}

		m.GetLogger().Debug(
			"metric updated",
			"type", mType,
			"name", name,
			"value", parseInt,
		)

	case models.Gauge:
		parseFloat, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("Value conversion error: %s", err), http.StatusBadRequest)
			return
		}

		_, err = m.metricsService.UpdateGauge(r.Context(), name, parseFloat)
		if err != nil {
			http.Error(w, fmt.Sprintf("Value update error: %s", err), http.StatusInternalServerError)
			return
		}

		m.GetLogger().Debug(
			"metric updated",
			"type", mType,
			"name", name,
			"value", parseFloat,
		)

	default:
		http.Error(w, "Invalid type", http.StatusBadRequest)
		return
	}
}
