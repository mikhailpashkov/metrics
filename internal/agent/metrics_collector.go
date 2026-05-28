package agent

import (
	"context"
	"log/slog"
	"os"
	"time"

	models "github.com/mikhailpashkov/metrics/internal/model"
	"github.com/mikhailpashkov/metrics/internal/service"
)

type MetricsPoller interface {
	GetMetrics() ([]*models.Metrics, error)
}

type MetricsReporter interface {
	SendMetrics(metrics []*models.Metrics) error
}

type MetricsCollectorParams struct {
	PollInterval   time.Duration
	ReportInterval time.Duration
	PollCallback   func()
	ReportCallback func()
}

type MetricsCollector struct {
	logger   *slog.Logger
	service  service.MetricsService
	pollers  []MetricsPoller
	reporter MetricsReporter
	params   *MetricsCollectorParams
}

func NewMetricsCollector(
	logger *slog.Logger,
	service service.MetricsService,
	pollers []MetricsPoller,
	reporter MetricsReporter,
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
	metricsToRecord := make(chan []*models.Metrics, 64)

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
			metricsToRecord <- accumulated
			err = m.service.DeleteAll(context.Background())
			if err != nil {
				slog.Error("Error deleting all metrics", "err", err)
				os.Exit(1)
			}
			go m.params.ReportCallback()
		}
	}()

	go func() {
		for metricsBatch := range metricsToRecord {
			err := m.reporter.SendMetrics(metricsBatch)
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
