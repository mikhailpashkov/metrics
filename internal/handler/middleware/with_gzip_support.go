package middleware

import (
	"compress/gzip"
	"io"
	"log/slog"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

type gzipHandler struct {
	http.Handler
	logger *slog.Logger
}

func (g gzipWriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}

func (h *gzipHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		h.logger.Debug("skip gzip encoding: not supported on client",
			"url", r.URL,
			"host", r.Host,
		)
		h.Handler.ServeHTTP(w, r)
		return
	}

	h.logger.Debug("gzip encoding is supported")

	gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
	if err != nil {
		h.logger.Error("Failed to create new gzip writer", "err", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer gz.Close()

	w.Header().Set("Content-Encoding", "gzip")
	h.Handler.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
}

func WithGZIPSupport(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return &gzipHandler{next, logger}
	}
}
