package models

const (
	MetricsUpdatedEvent = "metrics_updated"
	MetricsDeletedEvent = "metrics_deleted"
)

type Event struct {
	ID  string
	Key EventKey
}

type EventKey string

type EventSubscriber struct {
	Name     string
	Callback func(Event)
}
