package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"

	"github.com/mikhailpashkov/metrics/internal/handler"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

type gzipHandler struct {
	handler.MHandler
}

func (g gzipWriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}

func (h *gzipHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		h.GetLogger().Debug("skip gzip encoding: not supported on client",
			"url", r.URL,
			"host", r.Host,
		)
		h.MHandler.ServeHTTP(w, r)
		return
	}

	h.GetLogger().Debug("gzip encoding is supported")

	gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer gz.Close()

	w.Header().Set("Content-Encoding", "gzip")
	h.MHandler.ServeHTTP(gzipWriter{ResponseWriter: w, Writer: gz}, r)
}

func WithGZIPSupport(handler handler.MHandler) handler.MHandler {
	return &gzipHandler{MHandler: handler}
}
