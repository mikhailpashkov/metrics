package handler

import (
	"log/slog"
	"net/http"
)

type MHandler interface {
	GetLogger() *slog.Logger
	GetUrlPatterns() []string
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
