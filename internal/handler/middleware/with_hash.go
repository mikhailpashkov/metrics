package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"

	"github.com/mikhailpashkov/metrics/internal/utils"
)

type checkHashHandler struct {
	http.Handler
	key    string
	logger *slog.Logger
}

type writeHashHandler struct {
	http.Handler
	key    string
	logger *slog.Logger
}

type hashResponseWriter struct {
	header http.Header
	body   bytes.Buffer
	status int
}

func (w *hashResponseWriter) Header() http.Header {
	return w.header
}

func (w *hashResponseWriter) Write(b []byte) (int, error) {
	return w.body.Write(b)
}

func (w *hashResponseWriter) WriteHeader(status int) {
	w.status = status
}

func (h *checkHashHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	gotHash := r.Header.Get(utils.HashHeaderKey)
	if gotHash == "" {
		h.logger.Debug("skip hash check: empty hash header")
		h.Handler.ServeHTTP(w, r)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("failed to read body", "err", err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	r.Body = io.NopCloser(bytes.NewBuffer(body))

	expectedHash := utils.Hash(body, []byte(h.key))

	if !utils.HashEqual(gotHash, expectedHash) {
		h.logger.Warn("hash mismatch")
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	h.Handler.ServeHTTP(w, r)
}

func (h *writeHashHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rw := &hashResponseWriter{
		header: make(http.Header),
		body:   bytes.Buffer{},
		status: http.StatusOK,
	}

	h.Handler.ServeHTTP(rw, r)

	responseBody := rw.body.Bytes()
	calculatedHash := utils.Hash(responseBody, []byte(h.key))

	h.logger.Debug("calculated hash", "hash", calculatedHash)

	rw.Header().Set(
		utils.HashHeaderKey,
		calculatedHash,
	)

	for k, v := range rw.header {
		w.Header()[k] = v
	}

	w.WriteHeader(rw.status)
	_, e := w.Write(rw.body.Bytes())
	if e != nil {
		h.logger.Error("failed to write body", "err", e)
	}
}

func WithHASHCheck(logger *slog.Logger, key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return &checkHashHandler{
			Handler: next,
			key:     key,
			logger:  logger,
		}
	}
}

func WithHASHWrite(logger *slog.Logger, key string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return &writeHashHandler{
			Handler: next,
			key:     key,
			logger:  logger,
		}
	}
}
