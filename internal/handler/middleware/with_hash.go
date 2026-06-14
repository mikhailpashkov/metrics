package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log/slog"
	"net/http"
)

const hashHeaderKey = "HashSHA256"

func hash(body []byte, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))
}

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
	gotHash := r.Header.Get(hashHeaderKey)
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

	expectedHash := hash(body, []byte(h.key))

	if !hmac.Equal([]byte(gotHash), []byte(expectedHash)) {
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
	calculatedHash := hash(responseBody, []byte(h.key))

	h.logger.Debug("calculated hash", "hash", calculatedHash)

	rw.Header().Set(
		hashHeaderKey,
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
