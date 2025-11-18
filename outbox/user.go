package outbox

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// User represents our business entity
type User struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// UserCreatedEvent represents the event data
type UserCreatedEvent struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
}

// OutboxEvent represents a row in our outbox_events table
type OutboxEvent struct {
	ID            string          `json:"id"`
	AggregateType string          `json:"aggregate_type"`
	AggregateID   string          `json:"aggregate_id"`
	EventType     string          `json:"event_type"`
	Payload       json.RawMessage `json:"payload"`
	Metadata      json.RawMessage `json:"metadata"`
	CreatedAt     time.Time       `json:"created_at"`
	ProcessedAt   sql.NullTime    `json:"processed_at"` // Use sql.NullTime for nullable timestamps
}

type UserService struct {
	db *sql.DB
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{db: db}
}

// CreateUser creates a user and publishes a UserCreated event atomically
func (s *UserService) Create(ctx context.Context, name, email string) (*User, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, fmt.Errorf("begin transaction failed: %w", err)
	}
	defer handleRollback(tx)

	newUser := User{
		ID:        uuid.New().String(),
		Name:      name,
		Email:     email,
		CreatedAt: time.Now().UTC(),
	}
	userInsertStmt := `INSERT INTO users (id, name, email, created_at) VALUES ($1, $2, $3, $4)`
	_, err = tx.ExecContext(ctx, userInsertStmt, newUser.ID, newUser.Name, newUser.Email, newUser.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert user: %w", err)
	}

	log.Printf("✅ User %s created in database.", newUser.ID)
	// 2. Prepare the UserCreated event
	eventPayload, err := json.Marshal(UserCreatedEvent{UserID: newUser.ID, Email: newUser.Email})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event payload: %w", err)
	}
	outboxEvent := OutboxEvent{
		ID:            uuid.New().String(),
		AggregateType: "User",
		AggregateID:   newUser.ID,
		EventType:     "UserCreated",
		Payload:       eventPayload,
		CreatedAt:     time.Now().UTC(),
	} // 3. Insert the event into the outbox table within the same transaction
	outboxInsertStmt := `INSERT INTO outbox_events 
		(id, aggregate_type, aggregate_id, event_type, payload, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = tx.ExecContext(ctx, outboxInsertStmt,
		outboxEvent.ID, outboxEvent.AggregateType, outboxEvent.AggregateID,
		outboxEvent.EventType, outboxEvent.Payload, outboxEvent.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert outbox event: %w", err)
	}
	log.Printf("✉️ Outbox event %s stored in database.", outboxEvent.ID) // If both operations succeeded, commit the transaction
	if err = tx.Commit(); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}
	log.Println("🎉 Transaction committed successfully!")

	return &newUser, nil
}
