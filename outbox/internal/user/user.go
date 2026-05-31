package user

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/apex/log"
	"github.com/google/uuid"
	"gitlab.com/mshindle/tidbits/outbox/internal/outbox"
	"gitlab.com/mshindle/tidbits/outbox/internal/repository"
)

// User represents our business entity
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name" faker:"name"`
	Email     string    `json:"email" faker:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// CreatedEvent represents the event data
type CreatedEvent struct {
	UserID int64  `json:"user_id"`
	Email  string `json:"email"`
}

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
}

type Clock interface {
	Now() time.Time
}

type Service struct {
	repo       Repository
	or         outbox.Repository
	transactor repository.Transactor
	clock      Clock
}

func NewService(repo Repository, outboxRepo outbox.Repository, t repository.Transactor, clock Clock) *Service {
	return &Service{
		repo:       repo,
		or:         outboxRepo,
		transactor: t,
		clock:      clock,
	}
}

// Create creates a user and publishes a UserCreated event atomically
func (s *Service) Create(ctx context.Context, name, email string) (*User, error) {
	newUser := User{
		Name:      name,
		Email:     email,
		CreatedAt: s.clock.Now(),
	}

	err := s.transactor.Transact(ctx, func(txCtx context.Context) error {
		err := s.repo.CreateUser(txCtx, &newUser)
		if err != nil {
			return fmt.Errorf("failed to create user: %w", err)
		}
		log.WithFields(log.Fields{
			"id":         newUser.ID,
			"name":       newUser.Name,
			"email":      newUser.Email,
			"created_at": newUser.CreatedAt.Format(time.RFC3339),
		}).Infof("user created in database")

		// Prepare the UserCreated event
		eventPayload, err := json.Marshal(CreatedEvent{UserID: newUser.ID, Email: newUser.Email})
		if err != nil {
			return fmt.Errorf("failed to marshal event payload: %w", err)
		}
		outboxEvent := outbox.Event{
			ID:            uuid.New().String(),
			AggregateType: "User",
			AggregateID:   strconv.FormatInt(newUser.ID, 10),
			EventType:     "UserCreated",
			Payload:       eventPayload,
			CreatedAt:     time.Now().UTC(),
		}
		return s.or.CreateEvent(txCtx, &outboxEvent)
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &newUser, nil
}
