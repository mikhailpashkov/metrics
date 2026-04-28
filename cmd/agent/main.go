package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/mikhailpashkov/metrics/internal/agent"
	"github.com/mikhailpashkov/metrics/internal/agent/poller"
	"github.com/mikhailpashkov/metrics/internal/agent/reporter"
	"github.com/mikhailpashkov/metrics/internal/repository"
	"github.com/mikhailpashkov/metrics/internal/service"
	"github.com/mikhailpashkov/metrics/internal/utils"
)

func main() {

	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, opts))
	logger.Info("AGENT")

	var addr string
	var pollInterval int
	var reportInterval int

	utils.GetParams([]utils.Param{
		&utils.StringParam{
			EnvName:       "ADDRESS",
			FlagName:      "a",
			FlagUsage:     "backend server address",
			Default:       "localhost:8080",
			ValueConsumer: func(v string) { addr = v },
		},
		&utils.IntParam{
			EnvName:       "POLL_INTERVAL",
			FlagName:      "p",
			FlagUsage:     "poll interval in seconds",
			Default:       2,
			ValueConsumer: func(v int) { pollInterval = v },
		},
		&utils.IntParam{
			EnvName:       "REPORT_INTERVAL",
			FlagName:      "r",
			FlagUsage:     "report interval in seconds",
			Default:       10,
			ValueConsumer: func(v int) { reportInterval = v },
		},
	})

	logger.Info("params", "addr", addr)
	logger.Info("params", "pollInterval", pollInterval)
	logger.Info("params", "reportInterval", reportInterval)

	metricsRepository := repository.NewMetricsMemoryRepository()
	metricsService := service.NewMetricsService(metricsRepository)

	//consoleReporter := reporter.NewConsoleReporter()
	backendReporter := reporter.NewBackendReporter(addr, logger)

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
			PollInterval:   time.Duration(pollInterval) * time.Second,
			ReportInterval: time.Duration(reportInterval) * time.Second,
			PollCallback:   pollCountPoller.IncrementCount,
			ReportCallback: pollCountPoller.ResetCount,
		},
	)

	metricsCollector.Start()
}
