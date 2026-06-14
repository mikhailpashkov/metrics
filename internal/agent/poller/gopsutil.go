package poller

import (
	"strconv"

	models "github.com/mikhailpashkov/metrics/internal/model"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

type GoPsUtilPoller struct {
}

func NewGoPsUtilPoller() *GoPsUtilPoller { return &GoPsUtilPoller{} }

func (p *GoPsUtilPoller) GetMetrics() ([]*models.Metrics, error) {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	utilization, err := cpu.Percent(0, true)
	if err != nil {
		return nil, err
	}

	metrics := make([]*models.Metrics, 0)
	metrics = append(metrics, models.MetricsBuilderUnInt64("TotalMemory", memory.Total))
	metrics = append(metrics, models.MetricsBuilderUnInt64("FreeMemory", memory.Free))

	for i, percent := range utilization {
		metrics = append(metrics, models.MetricsBuilderFloat64("CPUutilization"+strconv.Itoa(i+1), percent))
	}

	return metrics, nil
}
