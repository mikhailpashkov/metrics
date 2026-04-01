package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mikhailpashkov/metrics/internal/service"
)

type GetListMetricsHandler struct {
	metricsService service.MetricsService
}

func NewGetListMetricsHandler(metricsService service.MetricsService) *GetListMetricsHandler {
	return &GetListMetricsHandler{
		metricsService: metricsService,
	}
}

func (m *GetListMetricsHandler) GetUrlPattern() string {
	return "/metrics"
}

func (m *GetListMetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
