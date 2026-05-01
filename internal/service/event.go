package service

import (
	"fmt"
	"log/slog"
	"sync"

	models "github.com/mikhailpashkov/metrics/internal/model"
)

type EventService interface {
	Notify(event *models.Event)
	Subscribe(key models.EventKey, subscriber *models.EventSubscriber)
}

const queueSize = 10

type InMemoryEventService struct {
	logger      *slog.Logger
	queue       chan models.Event
	subscribers map[models.EventKey][]*models.EventSubscriber
	mu          sync.RWMutex
}

func NewInMemoryEventService(logger *slog.Logger) *InMemoryEventService {
	return &InMemoryEventService{
		logger:      logger,
		queue:       make(chan models.Event, queueSize),
		subscribers: make(map[models.EventKey][]*models.EventSubscriber),
	}
}

func (es *InMemoryEventService) Notify(event *models.Event) {
	es.mu.Lock()
	defer es.mu.Unlock()

	es.logger.Debug("received event",
		"event_id", event.ID,
		"event_key", event.Key,
	)
	subs, ok := es.subscribers[event.Key]
	if !ok {
		es.logger.Debug("no subscribers found",
			"event_id", event.ID,
			"event_key", event.Key,
		)
		return
	}

	es.logger.Debug(fmt.Sprintf("found %d subscribers", len(subs)),
		"event_id", event.ID,
		"event_key", event.Key,
	)
	for _, sub := range es.subscribers[event.Key] {
		go func() {
			es.logger.Debug("start executing callback",
				"event_id", event.ID,
				"event_key", event.Key,
				"subscriber_name", sub.Name,
			)
			sub.Callback(*event)
			es.logger.Debug("callback execution finished",
				"event_id", event.ID,
				"event_key", event.Key,
				"subscriber_name", sub.Name,
			)
		}()
	}
}

func (es *InMemoryEventService) Subscribe(key models.EventKey, sub *models.EventSubscriber) {
	es.mu.Lock()
	defer es.mu.Unlock()

	es.subscribers[key] = append(es.subscribers[key], sub)
	es.logger.Debug("new subscriber registered",
		"subscription_key", key,
		"subscriber_name", sub.Name,
	)
}
