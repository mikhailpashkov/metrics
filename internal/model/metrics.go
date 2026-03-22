package models

const (
	Counter = "counter"
	Gauge   = "gauge"
)

type Metrics struct {
	ID    int64    `json:"id"`
	Type  string   `json:"type"`
	Name  string   `json:"name"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	TS    int64    `json:"ts"`
}
