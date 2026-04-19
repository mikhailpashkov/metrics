package reporter

import (
	"encoding/json"
	"fmt"

	models "github.com/mikhailpashkov/metrics/internal/model"
)

type ConsoleReporter struct{}

func NewConsoleReporter() MetricsReporter {
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
