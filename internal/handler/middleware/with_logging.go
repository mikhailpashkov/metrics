package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type fetchingInfoResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (h *fetchingInfoResponseWriter) WriteHeader(statusCode int) {
	h.ResponseWriter.WriteHeader(statusCode)
	h.statusCode = statusCode
}

func (h *fetchingInfoResponseWriter) Write(data []byte) (int, error) {
	bytes, err := h.ResponseWriter.Write(data)
	h.size += bytes
	return bytes, err
}

type loggingHandler struct {
	http.Handler
	logger *slog.Logger
}

func (h *loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	responseWriter := &fetchingInfoResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // логика с прописыванием 200 во внутренней реализации Write'а не сработает, поэтому по-дефолту ставим StatusOK
		size:           0,
	}

	h.Handler.ServeHTTP(responseWriter, r)

	duration := time.Since(startTime)

	h.logger.Info("request processed",
		"url", r.URL.String(),
		"method", r.Method,
		"duration", duration,
		"status", responseWriter.statusCode,
		"size", responseWriter.size,
	)
}

func WithLogging(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return &loggingHandler{next, logger}
	}
}
