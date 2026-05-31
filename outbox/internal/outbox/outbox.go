package outbox

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"
)

// Event represents a row in our outbox_events table
type Event struct {
	ID            string          `json:"id"`
	AggregateType string          `json:"aggregate_type"`
	AggregateID   string          `json:"aggregate_id"`
	EventType     string          `json:"event_type"`
	Payload       json.RawMessage `json:"payload"`
	Metadata      json.RawMessage `json:"metadata"`
	CreatedAt     time.Time       `json:"created_at"`
	ProcessedAt   sql.NullTime    `json:"processed_at"` // Use sql.NullTime for nullable timestamps
}

type Repository interface {
	//GetEventByID(id string) (*Event, error)
	CreateEvent(ctx context.Context, event *Event) error
	//UpdateEvent(event *Event) error
	//DeleteEvent(id string) error
}
