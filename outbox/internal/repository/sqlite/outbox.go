package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/apex/log"
	"gitlab.com/mshindle/tidbits/outbox/internal/outbox"
	"gitlab.com/mshindle/tidbits/outbox/internal/repository"
)

func (s *Store) CreateEvent(ctx context.Context, event *outbox.Event) error {
	var err error

	stmt := `INSERT INTO outbox_events
		(id, aggregate_type, aggregate_id, event_type, payload, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)`

	if tx, ok := ctx.Value(repository.SqliteTxKey).(*sql.Tx); ok {
		_, err = tx.ExecContext(ctx, stmt,
			event.ID, event.AggregateType, event.AggregateID,
			event.EventType, event.Payload, event.CreatedAt)
	} else {
		_, err = s.db.ExecContext(ctx, stmt,
			event.ID, event.AggregateType, event.AggregateID,
			event.EventType, event.Payload, event.CreatedAt)
	}
	if err != nil {
		return fmt.Errorf("failed to insert outbox event: %w", err)
	}
	log.WithField("id", event.ID).Info("outbox event stored in database")

	return nil
}
