package poller

import (
	"math/rand/v2"
	"time"

	models "github.com/mikhailpashkov/metrics/internal/model"
)

type RandomValuePoller struct{}

func NewRandomValuePoller() *RandomValuePoller { return &RandomValuePoller{} }

func (p *RandomValuePoller) GetMetrics() ([]*models.Metrics, error) {
	randomValue := rand.Float64()
	return []*models.Metrics{
		{
			ID:    -1,
			Type:  models.Gauge,
			Name:  "custom.RandomValue",
			Delta: nil,
			Value: &randomValue,
			TS:    time.Now().UnixMilli(),
		},
	}, nil
}
