package agent

import (
	"fmt"
	"time"

	models "github.com/mikhailpashkov/metrics/internal/model"
	"github.com/mikhailpashkov/metrics/internal/service"
)

type MetricsPoller interface {
	GetMetrics() ([]*models.Metrics, error)
}

type MetricsReporter interface {
	SendMetrics(metrics *models.Metrics) error
}

type MetricsCollectorParams struct {
	PollInterval   time.Duration
	ReportInterval time.Duration
}

type MetricsCollector struct {
	service  *service.MetricsService
	pollers  []MetricsPoller
	reporter MetricsReporter
	params   *MetricsCollectorParams
}

func NewMetricsCollector(
	service *service.MetricsService,
	pollers []MetricsPoller,
	reporter MetricsReporter,
	params *MetricsCollectorParams,
) *MetricsCollector {
	return &MetricsCollector{
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
			for _, poller := range m.pollers {
				metrics, err := poller.GetMetrics()
				if err != nil {
					fmt.Println("[ERR] Error polling metrics", err)
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
			_, err := m.service.UpdateMetrics(metric)
			if err != nil {
				fmt.Println("[ERR] Error updating metrics", metric, err)
				continue
			}
		}
	}()

	// report
	go func() {
		for {
			time.Sleep(m.params.ReportInterval)
			accumulated, err := m.service.GetAllAccumulated()
			if err != nil {
				fmt.Println("[ERR] Error getting all accumulated metrics", err)
			}
			for _, metric := range accumulated {
				metricsToRecord <- metric
			}
			err = m.service.DeleteAll()
			if err != nil {
				fmt.Println("[ERR] Error deleting all metrics", err)
				panic(err)
			}
		}
	}()

	go func() {
		for metric := range metricsToRecord {
			err := m.reporter.SendMetrics(metric)
			if err != nil {
				fmt.Println("[ERR] Error sending metrics to reporter", err)
				continue
			}
		}
	}()

	for {
	}
}
