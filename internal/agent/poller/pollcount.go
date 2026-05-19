package poller

import (
	"sync"
	"time"

	models "github.com/mikhailpashkov/metrics/internal/model"
)

type PollCountPoller struct {
	count int64
	mux   sync.Mutex
}

func NewPollCountPoller() *PollCountPoller {
	return &PollCountPoller{
		count: 0,
		mux:   sync.Mutex{},
	}
}

func (p *PollCountPoller) GetMetrics() ([]*models.Metrics, error) {
	p.mux.Lock()
	defer p.mux.Unlock()
	return []*models.Metrics{
		{
			ID:    -1,
			Type:  models.Counter,
			Name:  "PollCount",
			Delta: &p.count,
			Value: nil,
			TS:    time.Now().UnixMilli(),
		},
	}, nil
}

func (p *PollCountPoller) IncrementCount() {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.count++
}

func (p *PollCountPoller) ResetCount() {
	p.mux.Lock()
	defer p.mux.Unlock()
	p.count = 0
}
