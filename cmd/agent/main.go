package main

import (
	"flag"
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

	addr := flag.String("a", "localhost:8080", "backend server address")
	pollInterval := flag.Int("p", 2, "poll interval in seconds")
	reportInterval := flag.Int("r", 10, "report interval in seconds")

	flag.Parse()

	fmt.Println("addr", *addr)
	fmt.Println("pollInterval", *pollInterval)
	fmt.Println("reportInterval", *reportInterval)

	metricsRepository := repository.NewMetricsMemoryRepository()
	metricsService := service.NewMetricsService(metricsRepository)

	//consoleReporter := reporter.NewConsoleReporter()
	backendReporter := reporter.NewBackendReporter(*addr)

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
			PollInterval:   time.Duration(*pollInterval) * time.Second,
			ReportInterval: time.Duration(*reportInterval) * time.Second,
			PollCallback:   pollCountPoller.IncrementCount,
			ReportCallback: pollCountPoller.ResetCount,
		},
	)

	metricsCollector.Start()
}
