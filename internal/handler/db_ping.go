package handler

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

type DBPingHandler struct {
	logger *slog.Logger
	conn   *pgx.Conn
}

func NewDBPingHandler(logger *slog.Logger, conn *pgx.Conn) *DBPingHandler {
	return &DBPingHandler{
		logger: logger,
		conn:   conn,
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
	err := h.conn.Ping(dbPingTimeoutCtx)
	if err != nil {
		h.logger.Error("failed to ping DB", "err", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
}
