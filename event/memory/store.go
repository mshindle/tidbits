package memory

import (
	"sync"
	"time"

	"gitlab.com/mshindle/tidbits/event"
)

type EventStore struct {
	events []event.Event
	mu     sync.RWMutex
}

// NewEventStore creates a new event store instance
func NewEventStore() *EventStore {
	return &EventStore{
		events: make([]event.Event, 0, 2),
	}
}

// Save persists an event to the store
func (es *EventStore) Save(event event.Event) error {
	es.mu.Lock()
	defer es.mu.Unlock()
	es.events = append(es.events, event)
	return nil
}

// GetByType retrieves all events of a specific type
func (es *EventStore) GetByType(eventType string) ([]event.Event, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()

	var filtered []event.Event
	for _, e := range es.events {
		if e.Type == eventType {
			filtered = append(filtered, e)
		}
	}
	return filtered, nil
}

// GetAll retrieves all events
func (es *EventStore) GetAll() ([]event.Event, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()
	return append([]event.Event{}, es.events...), nil
}

// GetByTimeRange retrieves events within a time range
func (es *EventStore) GetByTimeRange(start, end time.Time) ([]event.Event, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()

	var filtered []event.Event
	for _, e := range es.events {
		if e.Timestamp.After(start) && e.Timestamp.Before(end) {
			filtered = append(filtered, e)
		}
	}
	return filtered, nil
}

// Count returns the total number of events
func (es *EventStore) Count() int {
	es.mu.RLock()
	defer es.mu.RUnlock()
	return len(es.events)
}

// GetLatest retrieves the N most recent events
func (es *EventStore) GetLatest(n int) ([]event.Event, error) {
	es.mu.RLock()
	defer es.mu.RUnlock()

	if n > len(es.events) {
		n = len(es.events)
	}

	start := len(es.events) - n
	return append([]event.Event{}, es.events[start:]...), nil
}
