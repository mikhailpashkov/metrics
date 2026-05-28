package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	models "github.com/mikhailpashkov/metrics/internal/model"
	"github.com/mikhailpashkov/metrics/internal/service"
)

func NewUpdateMetricsPathParamsHandlerFunc(logger *slog.Logger, metricsService service.MetricsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			_, err = metricsService.UpdateCounter(r.Context(), name, parseInt)
			if err != nil {
				logger.Error("Value update error", "err", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			logger.Debug(
				"metric updated",
				"type", mType,
				"name", name,
				"value", parseInt,
			)

		case models.Gauge:
			parseFloat, err := strconv.ParseFloat(valueStr, 64)
			if err != nil {
				logger.Error("Value conversion error", "err", err)
				http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
				return
			}

			_, err = metricsService.UpdateGauge(r.Context(), name, parseFloat)
			if err != nil {
				logger.Error("Value update error", "err", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			logger.Debug(
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
}
