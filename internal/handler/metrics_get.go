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

func NewGetMetricsHandlerFunc(logger *slog.Logger, metricsService service.MetricsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var request dto.GetMetricsRequest

		defer r.Body.Close()
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			logger.Debug("Error decoding body", "err", err)
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}

		if request.Type == "" {
			logger.Debug("empty type")
			http.Error(w, "Empty type", http.StatusBadRequest)
			return
		}

		if request.Type != models.Gauge && request.Type != models.Counter {
			logger.Debug("invalid type")
			http.Error(w, "Invalid type", http.StatusBadRequest)
			return
		}

		if request.Type == "" {
			logger.Debug("empty name")
			http.Error(w, "Empty name", http.StatusBadRequest)
			return
		}

		accumulated, err := metricsService.GetAllAccumulated(r.Context())
		if err != nil {
			logger.Error("Error getting all accumulated metrics", "err", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		var found []*models.Metrics
		for _, metrics := range accumulated {
			if metrics.Name == request.ID && metrics.Type == request.Type {
				found = append(found, metrics)
			}
		}

		if len(found) == 0 {
			logger.Debug("no metrics found")
			http.Error(w, "No metrics found", http.StatusNotFound)
			return
		}

		if len(found) > 1 {
			logger.Error("multiple metrics found")
			http.Error(w, "Multiple metrics found", http.StatusInternalServerError)
			return
		}

		response := mapper.MetricsToGetMetricsResponse(found[0])

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			logger.Error("error encoding response", "err", err)
			http.Error(w, "Cant encode json", http.StatusInternalServerError)
			return
		}
	}
}
