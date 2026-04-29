package poller

import (
	"runtime"
	"time"

	models "github.com/mikhailpashkov/metrics/internal/model"
)

type MemStatsPoller struct{}

func NewMemStatsPoller() *MemStatsPoller {
	return &MemStatsPoller{}
}
func (m MemStatsPoller) GetMetrics() ([]*models.Metrics, error) {
	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)

	metrics := make([]*models.Metrics, 0)
	metrics = append(metrics, m.metricsBuilderUnInt64("Alloc", ms.Alloc))
	metrics = append(metrics, m.metricsBuilderUnInt64("BuckHashSys", ms.BuckHashSys))
	metrics = append(metrics, m.metricsBuilderUnInt64("Frees", ms.Frees))
	metrics = append(metrics, m.metricsBuilderFloat64("GCCPUFraction", ms.GCCPUFraction))
	metrics = append(metrics, m.metricsBuilderUnInt64("GCSys", ms.GCSys))
	metrics = append(metrics, m.metricsBuilderUnInt64("HeapAlloc", ms.HeapAlloc))
	metrics = append(metrics, m.metricsBuilderUnInt64("HeapIdle", ms.HeapIdle))
	metrics = append(metrics, m.metricsBuilderUnInt64("HeapInuse", ms.HeapInuse))
	metrics = append(metrics, m.metricsBuilderUnInt64("HeapObjects", ms.HeapObjects))
	metrics = append(metrics, m.metricsBuilderUnInt64("HeapReleased", ms.HeapReleased))
	metrics = append(metrics, m.metricsBuilderUnInt64("HeapSys", ms.HeapSys))
	metrics = append(metrics, m.metricsBuilderUnInt64("LastGC", ms.LastGC))
	metrics = append(metrics, m.metricsBuilderUnInt64("Lookups", ms.Lookups))
	metrics = append(metrics, m.metricsBuilderUnInt64("MCacheInuse", ms.MCacheInuse))
	metrics = append(metrics, m.metricsBuilderUnInt64("MCacheSys", ms.MCacheSys))
	metrics = append(metrics, m.metricsBuilderUnInt64("MSpanInuse", ms.MSpanInuse))
	metrics = append(metrics, m.metricsBuilderUnInt64("MSpanSys", ms.MSpanSys))
	metrics = append(metrics, m.metricsBuilderUnInt64("Mallocs", ms.Mallocs))
	metrics = append(metrics, m.metricsBuilderUnInt64("NextGC", ms.NextGC))
	metrics = append(metrics, m.metricsBuilderUnInt64("NumForcedGC", uint64(ms.NumForcedGC)))
	metrics = append(metrics, m.metricsBuilderUnInt64("NumGC", uint64(ms.NumGC)))
	metrics = append(metrics, m.metricsBuilderUnInt64("OtherSys", ms.OtherSys))
	metrics = append(metrics, m.metricsBuilderUnInt64("PauseTotalNs", ms.PauseTotalNs))
	metrics = append(metrics, m.metricsBuilderUnInt64("StackInuse", ms.StackInuse))
	metrics = append(metrics, m.metricsBuilderUnInt64("StackSys", ms.StackSys))
	metrics = append(metrics, m.metricsBuilderUnInt64("Sys", ms.Sys))
	metrics = append(metrics, m.metricsBuilderUnInt64("TotalAlloc", ms.TotalAlloc))

	return metrics, nil
}

func (m MemStatsPoller) metricsBuilderUnInt64(name string, value uint64) *models.Metrics {
	return m.metricsBuilderFloat64(name, float64(value))
}

func (_ MemStatsPoller) metricsBuilderFloat64(name string, value float64) *models.Metrics {
	return &models.Metrics{
		ID:    -1,
		Type:  models.Gauge,
		Name:  name,
		Delta: nil,
		Value: &value,
		TS:    time.Now().UnixMilli(),
	}
}
