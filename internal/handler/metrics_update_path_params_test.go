package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mikhailpashkov/metrics/internal/repository"
	"github.com/mikhailpashkov/metrics/internal/service"
	"github.com/stretchr/testify/require"
)

func TestUpdateMetricsHandler_ServeHTTP(t *testing.T) {
	repo := repository.NewMetricsMemoryRepository()
	svc := service.NewMetricsService(repo, &repository.StubBackupRepository{})
	handler := NewUpdateMetricsPathParamsHandler(slog.Default(), svc)

	mux := http.NewServeMux()
	mux.Handle("/update/{type}/{name}/{value}", handler)

	tests := []struct {
		name       string
		path       string
		wantStatus int
	}{
		{
			name:       "Counter valid",
			path:       "/update/counter/requests_total/42",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Gauge valid",
			path:       "/update/gauge/cpu_usage/3.14",
			wantStatus: http.StatusOK,
		},
		{
			name:       "Invalid type",
			path:       "/update/unknown/foo/10",
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Missing value",
			path:       "/update/counter/foo/",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, tt.path, nil)
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			require.Equal(t, tt.wantStatus, rr.Result().StatusCode)

			if tt.wantStatus == http.StatusOK {
				metrics, err := repo.FindAll(context.Background())
				require.NoError(t, err)
				require.NotEmpty(t, metrics)

				marshal, err := json.Marshal(metrics)
				require.NoError(t, err)
				fmt.Println("SAVED METRICS:", string(marshal))

				err = repo.DeleteAll(context.Background())
				require.NoError(t, err)
			}
		})
	}
}
