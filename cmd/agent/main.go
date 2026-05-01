package main

import (
	"log/slog"
	"os"
	"time"

	"github.com/mikhailpashkov/metrics/internal/agent"
	"github.com/mikhailpashkov/metrics/internal/agent/poller"
	"github.com/mikhailpashkov/metrics/internal/agent/reporter"
	_const "github.com/mikhailpashkov/metrics/internal/config/const"
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
	var reportToLog bool

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
		&utils.BoolParam{
			EnvName:       "REPORT_TO_LOG",
			FlagName:      "report-to-log",
			FlagUsage:     "report-to-log - run without sending data to server. send data to logger",
			Default:       false,
			ValueConsumer: func(v bool) { reportToLog = v },
		},
	})

	logger.Info("params",
		"addr", addr,
		"pollInterval", pollInterval,
		"reportInterval", reportInterval,
		"reportToLog", reportToLog,
	)

	metricsRepository := repository.NewMetricsMemoryRepository()
	eventService := service.NewInMemoryEventService(logger.With(_const.LoggerNameKey, "service.InMemoryEventService"))
	metricsService := service.NewMetricsService(logger.With(_const.LoggerNameKey, "service.MetricsService"), metricsRepository, eventService)

	var metricsReporter reporter.MetricsReporter
	if reportToLog {
		metricsReporter = reporter.NewLogReporter(logger.With(_const.LoggerNameKey, "reporter.LogReporter"))
	} else {
		metricsReporter = reporter.NewBackendReporter(addr, logger.With(_const.LoggerNameKey, "reporter.BackendReporter"))
	}

	memStatsPoller := poller.NewMemStatsPoller()
	pollCountPoller := poller.NewPollCountPoller()
	randomValuePoller := poller.NewRandomValuePoller()

	metricsCollector := agent.NewMetricsCollector(
		logger.With(_const.LoggerNameKey, "agent.MetricsCollector"),
		metricsService,
		[]poller.MetricsPoller{
			memStatsPoller,
			pollCountPoller,
			randomValuePoller,
		},
		metricsReporter,
		&agent.MetricsCollectorParams{
			PollInterval:   time.Duration(pollInterval) * time.Second,
			ReportInterval: time.Duration(reportInterval) * time.Second,
			PollCallback:   pollCountPoller.IncrementCount,
			ReportCallback: pollCountPoller.ResetCount,
		},
	)

	metricsCollector.Start()
}
