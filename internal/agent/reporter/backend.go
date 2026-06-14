package reporter

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/mikhailpashkov/metrics/internal/dto"
	"github.com/mikhailpashkov/metrics/internal/mapper"
	models "github.com/mikhailpashkov/metrics/internal/model"
	"github.com/mikhailpashkov/metrics/internal/utils"
	"resty.dev/v3"
)

type BackendReporter struct {
	address string
	client  *resty.Client
	logger  *slog.Logger
}

type hashTransport struct {
	Key  string
	Base http.RoundTripper
}

func (t *hashTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.Body != nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body = io.NopCloser(bytes.NewReader(body))
		req.Header.Set(
			utils.HashHeaderKey,
			utils.Hash(body, []byte(t.Key)),
		)
	}
	return t.Base.RoundTrip(req)
}

func NewBackendReporter(address string, key string, logger *slog.Logger) *BackendReporter {
	client := resty.New()
	client.SetTimeout(5 * time.Second)
	if key != "" {
		client.SetTransport(&hashTransport{
			Key:  key,
			Base: http.DefaultTransport,
		})
	}
	return &BackendReporter{
		address: address,
		client:  client,
		logger:  logger,
	}
}

func (r *BackendReporter) SendMetrics(metrics []*models.Metrics) error {
	if len(metrics) == 0 {
		return nil
	}

	metricsDtos := make([]dto.MetricsDto, len(metrics))
	for i, m := range metrics {
		if m == nil {
			return fmt.Errorf("metrics cannot be nil")
		}
		if !models.IsValidMetrics(m) {
			return fmt.Errorf("invalid %s", m.Name)
		}
		metricsDtos[i] = mapper.MetricsToMetricsDto(m)
	}

	updateUrl := fmt.Sprintf("http://%s/updates",
		r.address,
	)

	requestBody := dto.UpdateMetricsBatchRequest(metricsDtos)

	request := r.client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept-Encoding", "gzip").
		SetBody(requestBody).
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second).
		SetRetryMaxWaitTime(5 * time.Second).
		AddRetryConditions(func(r *resty.Response, err error) bool {
			return err != nil || r.StatusCode() >= 500
		})

	resp, err := request.Post(updateUrl)
	if err != nil {
		return fmt.Errorf("update metrics failed: %w", err)
	}

	defer resp.Body.Close()
	if resp.StatusCode() != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read response body error: %w", err)
		}
		return fmt.Errorf("update metrics failed: unexpected status %d: %s", resp.StatusCode(), string(body))
	}

	r.logger.Debug("update metrics successfully", "count", len(metricsDtos))

	return nil
}
