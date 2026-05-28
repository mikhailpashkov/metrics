package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/mikhailpashkov/metrics/db/metrics"
)

func NewDBPingHandlerFunc(logger *slog.Logger, dbQueries *metrics.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		dbPingTimeoutCtx, cancelFunc := context.WithTimeout(r.Context(), 15*time.Second)
		defer cancelFunc()
		err := dbQueries.HealthCheck(dbPingTimeoutCtx)
		if err != nil {
			logger.Error("failed to healthcheck DB", "err", err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
	}

}
