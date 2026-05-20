package reporter

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	models "github.com/mikhailpashkov/metrics/internal/model"
	"github.com/stretchr/testify/require"
)

func ptrInt64(v int64) *int64       { return &v }
func ptrFloat64(v float64) *float64 { return &v }

func TestSendMetrics(t *testing.T) {
	tests := []struct {
		name           string
		metric         *models.Metrics
		serverResponse string
		serverStatus   int
		wantErr        bool
	}{
		{
			name: "Counter success",
			metric: &models.Metrics{
				Type:  models.Counter,
				Name:  "requests_total",
				Delta: ptrInt64(100),
			},
			serverResponse: "",
			serverStatus:   http.StatusOK,
			wantErr:        false,
		},
		{
			name: "Gauge success",
			metric: &models.Metrics{
				Type:  models.Gauge,
				Name:  "cpu_usage",
				Value: ptrFloat64(3.123456),
			},
			serverResponse: "",
			serverStatus:   http.StatusOK,
			wantErr:        false,
		},
		{
			name: "Server error",
			metric: &models.Metrics{
				Type:  models.Counter,
				Name:  "requests_total",
				Delta: ptrInt64(1),
			},
			serverResponse: "internal error",
			serverStatus:   http.StatusInternalServerError,
			wantErr:        true,
		},
		{
			name:    "Nil metric",
			metric:  nil,
			wantErr: true,
		},
		{
			name: "Counter with nil delta",
			metric: &models.Metrics{
				Type: models.Counter,
				Name: "missing_delta",
			},
			wantErr: true,
		},
		{
			name: "Gauge with nil value",
			metric: &models.Metrics{
				Type: models.Gauge,
				Name: "missing_value",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.serverStatus)
				_, err := fmt.Fprint(w, tt.serverResponse)
				require.NoError(t, err, "failed to write server response")
				require.Equal(t, http.MethodPost, r.Method)
				require.Equal(t, "application/json", r.Header.Get("Content-Type"))
				if tt.metric != nil {
					defer r.Body.Close()
					bodyBytes, err := io.ReadAll(r.Body)
					require.NoError(t, err)
					bodyString := string(bodyBytes)
					require.Contains(t, bodyString, tt.metric.Name, "body want to contain metric name")
				}
			}))
			defer ts.Close()

			br := NewBackendReporter(ts.Listener.Addr().String(), slog.Default())

			err := br.SendMetrics([]*models.Metrics{tt.metric})
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
