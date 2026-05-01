package agent

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/mikhailpashkov/metrics/internal/agent/poller"
	"github.com/mikhailpashkov/metrics/internal/agent/reporter"
	models "github.com/mikhailpashkov/metrics/internal/model"
	"github.com/mikhailpashkov/metrics/internal/service"
)

type MetricsCollectorParams struct {
	PollInterval   time.Duration
	ReportInterval time.Duration
	PollCallback   func()
	ReportCallback func()
}

type MetricsCollector struct {
	logger   *slog.Logger
	service  service.MetricsService
	pollers  []poller.MetricsPoller
	reporter reporter.MetricsReporter
	params   *MetricsCollectorParams
}

func NewMetricsCollector(
	logger *slog.Logger,
	service service.MetricsService,
	pollers []poller.MetricsPoller,
	reporter reporter.MetricsReporter,
	params *MetricsCollectorParams,
) *MetricsCollector {
	return &MetricsCollector{
		logger:   logger,
		service:  service,
		pollers:  pollers,
		reporter: reporter,
		params:   params,
	}
}

func (m *MetricsCollector) Start() {
	metricsToSave := make(chan *models.Metrics, 128)
	metricsToRecord := make(chan *models.Metrics, 128)

	// poll
	go func() {
		for {
			time.Sleep(m.params.PollInterval)
			go m.params.PollCallback()
			for _, metricsPoller := range m.pollers {
				metrics, err := metricsPoller.GetMetrics()
				if err != nil {
					m.logger.Error("Error polling metrics", "err", err)
					continue
				}
				for _, metric := range metrics {
					metricsToSave <- metric
				}
			}
		}
	}()

	// local save
	go func() {
		for metric := range metricsToSave {
			_, err := m.service.UpdateMetrics(context.Background(), metric)
			if err != nil {
				m.logger.Error("Error updating metrics", "err", err)
				continue
			}
		}
	}()

	// report
	go func() {
		for {
			time.Sleep(m.params.ReportInterval)
			accumulated, err := m.service.GetAllAccumulated(context.Background())
			if err != nil {
				m.logger.Error("Error getting all accumulated metrics", "err", err)
				continue
			}
			for _, metric := range accumulated {
				metricsToRecord <- metric
			}
			err = m.service.DeleteAll(context.Background())
			if err != nil {
				slog.Error("Error deleting all metrics", "err", err)
				os.Exit(1)
			}
			go m.params.ReportCallback()
		}
	}()

	go func() {
		for metric := range metricsToRecord {
			err := m.reporter.SendMetrics(metric)
			if err != nil {
				m.logger.Error("Error sending metrics to reporter", "err", err)
				continue
			}
		}
	}()

	for {
		time.Sleep(1 * time.Second)
	}
}
