package poller

import (
	"runtime"

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
	metrics = append(metrics, models.MetricsBuilderUnInt64("Alloc", ms.Alloc))
	metrics = append(metrics, models.MetricsBuilderUnInt64("BuckHashSys", ms.BuckHashSys))
	metrics = append(metrics, models.MetricsBuilderUnInt64("Frees", ms.Frees))
	metrics = append(metrics, models.MetricsBuilderFloat64("GCCPUFraction", ms.GCCPUFraction))
	metrics = append(metrics, models.MetricsBuilderUnInt64("GCSys", ms.GCSys))
	metrics = append(metrics, models.MetricsBuilderUnInt64("HeapAlloc", ms.HeapAlloc))
	metrics = append(metrics, models.MetricsBuilderUnInt64("HeapIdle", ms.HeapIdle))
	metrics = append(metrics, models.MetricsBuilderUnInt64("HeapInuse", ms.HeapInuse))
	metrics = append(metrics, models.MetricsBuilderUnInt64("HeapObjects", ms.HeapObjects))
	metrics = append(metrics, models.MetricsBuilderUnInt64("HeapReleased", ms.HeapReleased))
	metrics = append(metrics, models.MetricsBuilderUnInt64("HeapSys", ms.HeapSys))
	metrics = append(metrics, models.MetricsBuilderUnInt64("LastGC", ms.LastGC))
	metrics = append(metrics, models.MetricsBuilderUnInt64("Lookups", ms.Lookups))
	metrics = append(metrics, models.MetricsBuilderUnInt64("MCacheInuse", ms.MCacheInuse))
	metrics = append(metrics, models.MetricsBuilderUnInt64("MCacheSys", ms.MCacheSys))
	metrics = append(metrics, models.MetricsBuilderUnInt64("MSpanInuse", ms.MSpanInuse))
	metrics = append(metrics, models.MetricsBuilderUnInt64("MSpanSys", ms.MSpanSys))
	metrics = append(metrics, models.MetricsBuilderUnInt64("Mallocs", ms.Mallocs))
	metrics = append(metrics, models.MetricsBuilderUnInt64("NextGC", ms.NextGC))
	metrics = append(metrics, models.MetricsBuilderUnInt64("NumForcedGC", uint64(ms.NumForcedGC)))
	metrics = append(metrics, models.MetricsBuilderUnInt64("NumGC", uint64(ms.NumGC)))
	metrics = append(metrics, models.MetricsBuilderUnInt64("OtherSys", ms.OtherSys))
	metrics = append(metrics, models.MetricsBuilderUnInt64("PauseTotalNs", ms.PauseTotalNs))
	metrics = append(metrics, models.MetricsBuilderUnInt64("StackInuse", ms.StackInuse))
	metrics = append(metrics, models.MetricsBuilderUnInt64("StackSys", ms.StackSys))
	metrics = append(metrics, models.MetricsBuilderUnInt64("Sys", ms.Sys))
	metrics = append(metrics, models.MetricsBuilderUnInt64("TotalAlloc", ms.TotalAlloc))

	return metrics, nil
}
