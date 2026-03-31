package reporter

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	models "github.com/mikhailpashkov/metrics/internal/model"
)

type BackendReporter struct {
	address string
	client  *http.Client
}

func NewBackendReporter(address string) *BackendReporter {
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	return &BackendReporter{
		address: address,
		client:  client,
	}
}

func (r *BackendReporter) SendMetrics(metrics *models.Metrics) error {
	if metrics == nil {
		return fmt.Errorf("metrics cannot be nil")
	}

	var metricsValue string

	switch metrics.Type {
	case models.Counter:
		if metrics.Delta == nil {
			return fmt.Errorf("metrics delta is nil")
		}
		metricsValue = strconv.FormatInt(*metrics.Delta, 10)
	case models.Gauge:
		if metrics.Value == nil {
			return fmt.Errorf("metrics value is nil")
		}
		metricsValue = strconv.FormatFloat(*metrics.Value, 'f', -1, 64)
	default:
		return fmt.Errorf("unknown metric type: %s", metrics.Type)
	}

	updateUrl := fmt.Sprintf("http://%s/update/%s/%s/%s",
		r.address,
		metrics.Type,
		url.PathEscape(metrics.Name),
		url.PathEscape(metricsValue),
	)

	request, err := http.NewRequest(http.MethodPost, updateUrl, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	request.Header.Set("Content-Type", "text/plain")

	resp, err := r.client.Do(request)
	if err != nil {
		return fmt.Errorf("update metrics failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read response body error: %w", err)
		}
		return fmt.Errorf("update metrics failed: unexpected status %d: %s", resp.StatusCode, string(body))
	}

	fmt.Println("BackendReporter - sent update for", metrics.Type, metrics.Name, metricsValue)

	return nil
}
