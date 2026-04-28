package handler

import (
	"log/slog"
	"net/http"
)

type MHandler interface {
	GetLogger() *slog.Logger
	GetUrlPattern() string
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
