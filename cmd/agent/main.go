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

	//consoleReporter := reporter.NewConsoleReporter()
	backendReporter := reporter.NewBackendReporter("localhost:8080")

	memStatsPoller := poller.NewMemStatsPoller()
	pollCountPoller := poller.NewPollCountPoller()
	randomValuePoller := poller.NewRandomValuePoller()

	metricsCollector := agent.NewMetricsCollector(
		metricsService,
		[]poller.MetricsPoller{
			memStatsPoller,
			pollCountPoller,
			randomValuePoller,
		},
		backendReporter,
		&agent.MetricsCollectorParams{
			PollInterval:   1 * time.Second,
			ReportInterval: 10 * time.Second,
			PollCallback:   pollCountPoller.IncrementCount,
			ReportCallback: pollCountPoller.ResetCount,
		},
	)

	metricsCollector.Start()
}
