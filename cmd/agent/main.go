package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mikhailpashkov/metrics/internal/agent"
	models "github.com/mikhailpashkov/metrics/internal/model"
	"github.com/mikhailpashkov/metrics/internal/repository"
	"github.com/mikhailpashkov/metrics/internal/service"
)
import "runtime"

func main() {
	fmt.Println("AGENT")

	metricsRepository := repository.NewMetricsMemoryRepository()
	metricsService := service.NewMetricsService(metricsRepository)
	consoleReporter := NewConsoleReporter()

	metricsCollector := agent.NewMetricsCollector(
		metricsService,
		[]agent.MetricsPoller{MemStatsPoller{}},
		consoleReporter,
		&agent.MetricsCollectorParams{
			PollInterval:   1 * time.Second,
			ReportInterval: 10 * time.Second,
		},
	)

	metricsCollector.Start()

}

type MemStatsPoller struct{}

func (m MemStatsPoller) GetMetrics() ([]*models.Metrics, error) {
	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)

	metrics := make([]*models.Metrics, 0)
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.Alloc", ms.Alloc))
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.BuckHashSys", ms.BuckHashSys))
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.Frees", ms.Frees))
	metrics = append(metrics, m.metricsBuilderFloat64("ms.GCCPUFraction", ms.GCCPUFraction))
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.GCSys", ms.GCSys))
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.HeapAlloc", ms.HeapAlloc))
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.HeapIdle", ms.HeapIdle))
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.HeapInuse", ms.HeapInuse))
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.HeapObjects", ms.HeapObjects))
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.HeapReleased", ms.HeapReleased))
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.HeapSys", ms.HeapSys))
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.LastGC", ms.LastGC))
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.Lookups", ms.Lookups))
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.MCacheInuse", ms.MCacheInuse))
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.MCacheSys", ms.MCacheSys))
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.MSpanInuse", ms.MSpanInuse))
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.MSpanSys", ms.MSpanSys))
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.Mallocs", ms.Mallocs))
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.NextGC", ms.NextGC))
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.NumForcedGC", uint64(ms.NumForcedGC)))
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.NumGC", uint64(ms.NumGC)))
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.OtherSys", ms.OtherSys))
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.PauseTotalNs", ms.PauseTotalNs))
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.StackInuse", ms.StackInuse))
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.StackSys", ms.StackSys))
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.Sys", ms.Sys))
	metrics = append(metrics, m.metricsBuilderUnInt64("ms.TotalAlloc", ms.TotalAlloc))

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

type ConsoleReporter struct{}

func NewConsoleReporter() agent.MetricsReporter {
	return &ConsoleReporter{}
}
func (c ConsoleReporter) SendMetrics(metrics *models.Metrics) error {
	marshal, err := json.Marshal(metrics)
	if err != nil {
		return err
	}
	fmt.Println(string(marshal))
	return nil
}
