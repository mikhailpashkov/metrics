package handler

import (
	"net/http"
	"strconv"

	models "github.com/mikhailpashkov/metrics/internal/model"
	"github.com/mikhailpashkov/metrics/internal/service"
	"go.uber.org/zap"
)

type GetMetricsHandler struct {
	logger         *zap.Logger
	metricsService service.MetricsService
}

func NewGetMetricsHandler(logger *zap.Logger, metricsService service.MetricsService) *GetMetricsHandler {
	return &GetMetricsHandler{
		logger:         logger,
		metricsService: metricsService,
	}
}

func (m *GetMetricsHandler) GetLogger() *zap.Logger { return m.logger }

func (m *GetMetricsHandler) GetUrlPattern() string {
	return "/value/{type}/{name}"
}

func (m *GetMetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	mType := r.PathValue("type")
	name := r.PathValue("name")

	if mType == "" {
		http.Error(w, "Empty type", http.StatusBadRequest)
		return
	}

	if mType != models.Gauge && mType != models.Counter {
		http.Error(w, "Invalid type", http.StatusBadRequest)
		return
	}

	if name == "" {
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
		if metrics.Name == name && metrics.Type == mType {
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

	result := found[0]

	switch result.Type {
	case models.Gauge:
		_, err := w.Write([]byte(strconv.FormatFloat(*result.Value, 'f', -1, 64)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case models.Counter:
		_, err := w.Write([]byte(strconv.FormatInt(*result.Delta, 10)))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

}
