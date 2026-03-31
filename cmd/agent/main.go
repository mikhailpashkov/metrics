package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/mikhailpashkov/metrics/internal/agent"
	"github.com/mikhailpashkov/metrics/internal/agent/poller"
	models "github.com/mikhailpashkov/metrics/internal/model"
	"github.com/mikhailpashkov/metrics/internal/repository"
	"github.com/mikhailpashkov/metrics/internal/service"
)

func main() {
	fmt.Println("AGENT")

	metricsRepository := repository.NewMetricsMemoryRepository()
	metricsService := service.NewMetricsService(metricsRepository)
	consoleReporter := NewConsoleReporter()

	memStatsPoller := poller.NewMemStatsPoller()
	pollCountPoller := poller.NewPollCountPoller()
	randomValuePoller := poller.NewRandomValuePoller()

	metricsCollector := agent.NewMetricsCollector(
		metricsService,
		[]agent.MetricsPoller{
			memStatsPoller,
			pollCountPoller,
			randomValuePoller,
		},
		consoleReporter,
		&agent.MetricsCollectorParams{
			PollInterval:   1 * time.Second,
			ReportInterval: 10 * time.Second,
			PollCallback:   pollCountPoller.IncrementCount,
		},
	)

	metricsCollector.Start()
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
