package middleware

import (
	"net/http"
	"time"

	"github.com/mikhailpashkov/metrics/internal/handler"
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
	handler.MHandler
}

func (h *loggingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	responseWriter := &fetchingInfoResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // логика с прописыванием 200 во внутренней реализации Write'а не сработает, поэтому по-дефолту ставим StatusOK
		size:           0,
	}

	h.MHandler.ServeHTTP(responseWriter, r)

	duration := time.Since(startTime)

	h.MHandler.GetLogger().Info("request processed",
		"url", r.URL.String(),
		"method", r.Method,
		"duration", duration,
		"status", responseWriter.statusCode,
		"size", responseWriter.size,
	)
}

func WithLogging(handler handler.MHandler) handler.MHandler {
	return &loggingHandler{MHandler: handler}
}
