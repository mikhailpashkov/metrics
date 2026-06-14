package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"log/slog"
	"net/http"
	"slices"
)

const hashHeaderKey = "HashSHA256"

var supportedMethods = []string{http.MethodPost, http.MethodPut, http.MethodPatch}

type hashHandler struct {
	http.Handler
	key    string
	logger *slog.Logger
}

type hashResponseWriter struct {
	http.ResponseWriter
	body bytes.Buffer
}

func (w *hashResponseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func hash(body []byte, key []byte) string {
	h := hmac.New(sha256.New, key)
	h.Write(body)
	return hex.EncodeToString(h.Sum(nil))
}

func (h *hashHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if !slices.Contains(supportedMethods, r.Method) {
		h.logger.Debug("skip hash check: not supported method")
		h.Handler.ServeHTTP(w, r)
		return
	}

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

	rw := &hashResponseWriter{
		ResponseWriter: w,
	}

	h.Handler.ServeHTTP(rw, r)

	responseBody := rw.body.Bytes()

	w.Header().Set(
		hashHeaderKey,
		hash(responseBody, []byte(h.key)),
	)
}

func WithHASH(logger *slog.Logger, key string) func(next http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {
		return &hashHandler{
			Handler: next,
			key:     key,
			logger:  logger,
		}
	}
}
