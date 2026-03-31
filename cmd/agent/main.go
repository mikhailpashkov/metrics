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
