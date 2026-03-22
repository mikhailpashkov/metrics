package handler

import (
	"fmt"
	"net/http"
	"strconv"

	models "github.com/mikhailpashkov/metrics/internal/model"
	"github.com/mikhailpashkov/metrics/internal/service"
)

type Metrics struct {
	metricsService *service.Metrics
}

func NewMetrics(metricsService *service.Metrics) *Metrics {
	return &Metrics{
		metricsService: metricsService,
	}
}

func (m *Metrics) GetUrlPattern() string {
	return "/update/{type}/{name}/{value}"
}

func (m *Metrics) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	if mType != models.Gauge && mType != models.Counter {
		http.Error(w, "Invalid type", http.StatusBadRequest)
		return
	}

	if name == "" {
		http.Error(w, "Empty name", http.StatusBadRequest)
		return
	}

	if mType == models.Counter {
		parseInt, err := strconv.ParseInt(valueStr, 10, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("Value conversion error: %s", err), http.StatusBadRequest)
			return
		}

		_, err = m.metricsService.UpdateCounter(name, parseInt)
		if err != nil {
			http.Error(w, fmt.Sprintf("Value update error: %s", err), http.StatusInternalServerError)
			return
		}
	}

	if mType == models.Gauge {
		parseFloat, err := strconv.ParseFloat(valueStr, 64)
		if err != nil {
			http.Error(w, fmt.Sprintf("Value conversion error: %s", err), http.StatusBadRequest)
			return
		}

		_, err = m.metricsService.UpdateGauge(name, parseFloat)
		if err != nil {
			http.Error(w, fmt.Sprintf("Value update error: %s", err), http.StatusInternalServerError)
			return
		}
	}
}
