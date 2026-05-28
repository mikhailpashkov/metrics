package dto

type GetMetricsRequest struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type MetricsDto struct {
	ID    string   `json:"id"`
	Type  string   `json:"type"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
}

type GetMetricsResponse MetricsDto

type UpdateMetricsRequest MetricsDto
type UpdateMetricsBatchRequest []MetricsDto
