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

type GetMetricsHandler struct {
	logger         *slog.Logger
	metricsService service.MetricsService
}

func NewGetMetricsHandler(logger *slog.Logger, metricsService service.MetricsService) *GetMetricsHandler {
	return &GetMetricsHandler{
		logger:         logger,
		metricsService: metricsService,
	}
}

func (m *GetMetricsHandler) GetLogger() *slog.Logger { return m.logger }

func (m *GetMetricsHandler) GetUrlPattern() string {
	return "/value"
}

func (m *GetMetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	var request dto.GetMetricsRequest

	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "invalid request: "+err.Error(), http.StatusBadRequest)
		return
	}

	if request.Type == "" {
		http.Error(w, "Empty type", http.StatusBadRequest)
		return
	}

	if request.Type != models.Gauge && request.Type != models.Counter {
		http.Error(w, "Invalid type", http.StatusBadRequest)
		return
	}

	if request.Type == "" {
		http.Error(w, "Empty name", http.StatusBadRequest)
		return
	}

	accumulated, err := m.metricsService.GetAllAccumulated(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var found []*models.Metrics
	for _, metrics := range accumulated {
		if metrics.Name == request.ID && metrics.Type == request.Type {
			found = append(found, metrics)
		}
	}

	if len(found) == 0 {
		http.Error(w, "No metrics found", http.StatusNotFound)
		return
	}

	if len(found) > 1 {
		http.Error(w, "Multiple metrics found", http.StatusInternalServerError)
		return
	}

	response, err := mapper.MetricsToGetMetricsResponse(found[0])
	if err != nil {
		http.Error(w, "Can't map MetricsToGetMetricsResponse: "+err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Cant encode json: "+err.Error(), http.StatusInternalServerError)
		return
	}
}
