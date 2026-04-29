package dto

import "encoding/json"

type GetMetricsRequest struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type GetMetricsResponse struct {
	ID    string      `json:"id"`
	Type  string      `json:"type"`
	Value json.Number `json:"value"`
}

type UpdateMetricsRequest struct {
	ID    string      `json:"id"`
	Type  string      `json:"type"`
	Value json.Number `json:"value"`
}
