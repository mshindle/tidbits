package event

import (
	"sync"
	"time"

	"github.com/apex/log"
	"gitlab.com/mshindle/tidbits/errorh"
)

type Event struct {
	ID        string
	Type      string
	Payload   interface{}
	Timestamp time.Time
}

// Handler is a function that processes events
type Handler func(Event) error

// Bus manages event subscriptions and publishing (Pub/Sub pattern)
type Bus struct {
	handlers map[string][]Handler
	mu       sync.RWMutex
}

// NewBus creates a new event bus instance
func NewBus() *Bus {
	return &Bus{
		handlers: make(map[string][]Handler),
	}
}

func (b *Bus) Subscribe(topic string, handler Handler) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.handlers[topic] = append(b.handlers[topic], handler)
}

func (b *Bus) Publish(event Event) error {
	b.mu.RLock()
	handlers := make([]Handler, 0, 5)
	handlers = append(handlers, b.handlers[event.Type]...)
	handlers = append(handlers, b.handlers["*"]...)
	b.mu.RUnlock()

	errs := errorh.NewErrorCollection()
	for _, handler := range handlers {
		if err := handler(event); err != nil {
			log.WithError(err).WithField("event_type", event.Type).Error("handler error for event")
			errs.Add(err)
		}
		// Continue processing other handlers even if one fails
	}

	if errs.HasErrors() {
		return errs
	}
	return nil
}

func (b *Bus) PublishAsync(event Event) {
	go func() {
		if err := b.Publish(event); err != nil {
			log.WithError(err).Error("async publish error")
		}
	}()
}

// Unsubscribe removes all handlers for a specific event type
func (b *Bus) Unsubscribe(eventType string) {
	b.mu.Lock()
	defer b.mu.Unlock()
	delete(b.handlers, eventType)
}

// GetSubscriberCount returns the number of subscribers for an event type
func (b *Bus) GetSubscriberCount(eventType string) int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.handlers[eventType])
}
