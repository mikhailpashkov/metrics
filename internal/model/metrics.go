package models

import "time"

const (
	Counter = "counter"
	Gauge   = "gauge"
)

const MetricsNewID = -1

type Metrics struct {
	ID    int64    `json:"id"`
	Type  string   `json:"type"`
	Name  string   `json:"name"`
	Delta *int64   `json:"delta,omitempty"`
	Value *float64 `json:"value,omitempty"`
	TS    int64    `json:"ts"`
}

func MetricsBuilderUnInt64(name string, value uint64) *Metrics {
	return MetricsBuilderFloat64(name, float64(value))
}

func MetricsBuilderFloat64(name string, value float64) *Metrics {
	return &Metrics{
		ID:    MetricsNewID,
		Type:  Gauge,
		Name:  name,
		Delta: nil,
		Value: &value,
		TS:    time.Now().UnixMilli(),
	}
}
