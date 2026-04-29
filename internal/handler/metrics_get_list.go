package handler

import (
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"sort"

	"github.com/mikhailpashkov/metrics/internal/service"
)

type GetListMetricsHandler struct {
	logger         *slog.Logger
	metricsService service.MetricsService
}

const htmlTemplate = `<html>
<head>
  <meta charset="utf-8">
  <title>All Metrics</title>
</head>
<body>
  <h1>All Metrics</h1>
<table>
<thead>
<tr>
<td>Type</td>
<td>Name</td>
<td>Value</td>
</tr>
</thead>
<tbody>
  {{range $key, $value := .}}
	<tr>
<td>{{$value.Type}}</td>
<td>{{$value.Name}}</td>
{{ if eq $value.Type "gauge" }}
<td>{{$value.Value}}</td>
{{ else if eq $value.Type "counter" }}
<td>{{$value.Delta}}</td>
{{ end }}
</tr>
  {{end}}
</tbody>
</table>
</body>
</html>
`

func NewGetListMetricsHandler(logger *slog.Logger, metricsService service.MetricsService) *GetListMetricsHandler {
	return &GetListMetricsHandler{
		logger:         logger,
		metricsService: metricsService,
	}
}

func (m *GetListMetricsHandler) GetLogger() *slog.Logger { return m.logger }

func (m *GetListMetricsHandler) GetUrlPatterns() []string {
	return []string{"/"}
}

func (m *GetListMetricsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	accumulated, err := m.metricsService.GetAllAccumulated(r.Context())
	if err != nil {
		http.Error(w, fmt.Sprintf("GetAllAccumulated error: %s", err), http.StatusInternalServerError)
		return
	}
	sort.Slice(accumulated, func(i, j int) bool {
		return accumulated[i].Name < accumulated[j].Name
	})

	w.Header().Set("Content-Type", "text/html")

	tmpl, err := template.New("metrics").Parse(htmlTemplate)
	if err != nil {
		http.Error(w, fmt.Sprintf("template parse error: %s", err), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, accumulated)
	if err != nil {
		http.Error(w, fmt.Sprintf("template execute error: %s", err), http.StatusInternalServerError)
		return
	}
}
