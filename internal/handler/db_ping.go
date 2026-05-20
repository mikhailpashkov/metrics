package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/mikhailpashkov/metrics/db/metrics"
)

type DBPingHandler struct {
	logger    *slog.Logger
	dbQueries *metrics.Queries
}

func NewDBPingHandler(logger *slog.Logger, dbQueries *metrics.Queries) *DBPingHandler {
	return &DBPingHandler{
		logger:    logger,
		dbQueries: dbQueries,
	}
}

func (h *DBPingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		h.logger.Debug("Method not allowed")
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	dbPingTimeoutCtx, cancelFunc := context.WithTimeout(r.Context(), 15*time.Second)
	defer cancelFunc()
	err := h.dbQueries.HealthCheck(dbPingTimeoutCtx)
	if err != nil {
		h.logger.Error("failed to healthcheck DB", "err", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
