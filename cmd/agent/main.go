package main

import (
	"fmt"
	"time"

	"github.com/mikhailpashkov/metrics/internal/agent"
	"github.com/mikhailpashkov/metrics/internal/agent/poller"
	"github.com/mikhailpashkov/metrics/internal/agent/reporter"
	"github.com/mikhailpashkov/metrics/internal/repository"
	"github.com/mikhailpashkov/metrics/internal/service"
)

func main() {
	fmt.Println("AGENT")

	metricsRepository := repository.NewMetricsMemoryRepository()
	metricsService := service.NewMetricsService(metricsRepository)
	consoleReporter := reporter.NewConsoleReporter()

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
