package user

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gitlab.com/mshindle/tidbits/outbox/internal/outbox"
)

// MockRepository is a mock implementation of the Repository interface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateUser(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

// MockOutboxRepository is a mock implementation of the outbox.Repository interface
type MockOutboxRepository struct {
	mock.Mock
}

func (m *MockOutboxRepository) CreateEvent(ctx context.Context, event *outbox.Event) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

// MockClock is a mock implementation of the Clock interface
type MockClock struct {
	mock.Mock
}

func (m *MockClock) Now() time.Time {
	args := m.Called()
	return args.Get(0).(time.Time)
}

// MockTransactor is a mock implementation of the repository.Transactor interface
type MockTransactor struct {
	mock.Mock
}

func (m *MockTransactor) Transact(ctx context.Context, fn func(ctx context.Context) error) error {
	args := m.Called(ctx, fn)
	// Execute the function passed to Transact with the context it expects
	if fnErr := fn(ctx); fnErr != nil {
		return fnErr
	}
	return args.Error(0)
}

func TestNewService(t *testing.T) {
	repo := new(MockRepository)
	outboxRepo := new(MockOutboxRepository)
	transactor := new(MockTransactor)
	clock := new(MockClock)

	s := NewService(repo, outboxRepo, transactor, clock)

	assert.NotNil(t, s)
	assert.Equal(t, repo, s.repo)
	assert.Equal(t, outboxRepo, s.or)
	assert.Equal(t, transactor, s.transactor)
	assert.Equal(t, clock, s.clock)
}

func TestService_Create(t *testing.T) {
	ctx := t.Context()
	name := "John Doe"
	email := "john@example.com"
	now := time.Now().UTC()

	t.Run("success", func(t *testing.T) {
		repo := new(MockRepository)
		outboxRepo := new(MockOutboxRepository)
		transactor := new(MockTransactor)
		clock := new(MockClock)
		s := NewService(repo, outboxRepo, transactor, clock)

		transactor.On("Transact", ctx, mock.Anything).Return(nil)
		clock.On("Now").Return(now)
		repo.On("CreateUser", ctx, mock.MatchedBy(func(u *User) bool {
			u.ID = 123 // simulate DB setting the ID
			return u.Name == name && u.Email == email && u.CreatedAt.Equal(now)
		})).Return(nil)

		outboxRepo.On("CreateEvent", ctx, mock.MatchedBy(func(e *outbox.Event) bool {
			return e.AggregateType == "User" &&
				e.AggregateID == "123" &&
				e.EventType == "UserCreated" &&
				len(e.Payload) > 0
		})).Return(nil)

		user, err := s.Create(ctx, name, email)

		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, int64(123), user.ID)
		assert.Equal(t, name, user.Name)
		assert.Equal(t, email, user.Email)
		assert.Equal(t, now, user.CreatedAt)

		clock.AssertExpectations(t)
		repo.AssertExpectations(t)
		outboxRepo.AssertExpectations(t)
		transactor.AssertExpectations(t)
	})

	t.Run("repo error", func(t *testing.T) {
		repo := new(MockRepository)
		outboxRepo := new(MockOutboxRepository)
		transactor := new(MockTransactor)
		clock := new(MockClock)
		s := NewService(repo, outboxRepo, transactor, clock)

		transactor.On("Transact", ctx, mock.Anything).Return(nil)
		clock.On("Now").Return(now)
		repo.On("CreateUser", ctx, mock.Anything).Return(errors.New("db error"))

		user, err := s.Create(ctx, name, email)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create user: db error")
		assert.Nil(t, user)

		clock.AssertExpectations(t)
		repo.AssertExpectations(t)
		outboxRepo.AssertExpectations(t)
		transactor.AssertExpectations(t)
	})

	t.Run("outbox error", func(t *testing.T) {
		repo := new(MockRepository)
		outboxRepo := new(MockOutboxRepository)
		transactor := new(MockTransactor)
		clock := new(MockClock)
		s := NewService(repo, outboxRepo, transactor, clock)

		transactor.On("Transact", ctx, mock.Anything).Return(nil)
		clock.On("Now").Return(now)
		repo.On("CreateUser", ctx, mock.Anything).Return(nil)
		outboxRepo.On("CreateEvent", ctx, mock.Anything).Return(errors.New("outbox error"))

		user, err := s.Create(ctx, name, email)

		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create user: outbox error")
		assert.Nil(t, user)

		clock.AssertExpectations(t)
		repo.AssertExpectations(t)
		outboxRepo.AssertExpectations(t)
		transactor.AssertExpectations(t)
	})
}
