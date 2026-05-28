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

func NewUpdateMetricsHandlerFunc(logger *slog.Logger, metricsService service.MetricsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var request dto.UpdateMetricsRequest

		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			logger.Debug("Error decoding body", "err", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if request.Type == "" {
			logger.Debug("Invalid request: Empty type")
			http.Error(w, "Invalid request: Empty type", http.StatusBadRequest)
			return
		}

		if request.ID == "" {
			logger.Debug("Invalid request: Empty name")
			http.Error(w, "Invalid request: Empty name", http.StatusBadRequest)
			return
		}

		metrics := mapper.MetricsFromUpdateMetricsRequest(request)

		isValid := models.IsValidMetrics(metrics)
		if !isValid {
			logger.Debug("Metric type doesnt match its content")
			http.Error(w, "Metric type doesnt match its content", http.StatusBadRequest)
			return
		}

		_, err = metricsService.UpdateMetrics(r.Context(), metrics)
		if err != nil {
			logger.Error("Value update error", "err", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		_, err = w.Write([]byte("{}"))
		if err != nil {
			logger.Error("failed to write response", "err", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

	}
}
