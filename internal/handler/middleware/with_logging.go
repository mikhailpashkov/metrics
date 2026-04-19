package middleware

import (
	"net/http"
	"time"

	"github.com/mikhailpashkov/metrics/internal/handler"
	"go.uber.org/zap"
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
		statusCode:     http.StatusOK, // логика с прописыванием 200 во внутренней реализации Write'а не сработает, поэтому по-дефолту StatusOK
		size:           0,
	}

	h.MHandler.ServeHTTP(responseWriter, r)

	duration := time.Since(startTime)

	h.MHandler.GetLogger().Info("Request processed",
		zap.String("url", r.URL.String()),
		zap.String("method", r.Method),
		zap.Duration("duration", duration),
		zap.Int("status", responseWriter.statusCode),
		zap.Int("size", responseWriter.size),
	)
}

func WithLogging(handler handler.MHandler) handler.MHandler {
	return &loggingHandler{MHandler: handler}
}
