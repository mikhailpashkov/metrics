package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mikhailpashkov/metrics/internal/service"
)

type GetMetricsHandler struct {
	metricsService *service.MetricsService
}

func NewGetMetricsHandler(metricsService *service.MetricsService) *GetMetricsHandler {
	return &GetMetricsHandler{
		metricsService: metricsService,
	}
}

func (m *GetMetricsHandler) GetUrlPattern() string {
	return "/metrics"
}

func (m *GetMetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	allRecords, err := m.metricsService.GetAllRecords()
	if err != nil {
		http.Error(w, fmt.Sprintf("GetAllRecords error: %s", err), http.StatusInternalServerError)
		return
	}

	bytes, err := json.Marshal(allRecords)
	if err != nil {
		http.Error(w, fmt.Sprintf("json.Marshal error: %s", err), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(bytes)
}
