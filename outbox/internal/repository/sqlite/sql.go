package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/apex/log"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) InitializeDB() error {
	var err error

	_, err = s.db.Exec(CreateUsersTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create users table: %w", err)
	}
	log.Info("users table checked/created")

	_, err = s.db.Exec(CreateOutboxEventsTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create outbox_events table: %w", err)
	}
	log.Info("outbox events table checked/created")

	return nil
}

// CreateUsersTableSQL contains the SQL statement to create a users table
// that matches the User struct fields and types.
//
// Notes:
//   - ID uses auto-incrementing integer as the PRIMARY KEY.
//   - Name and Email are TEXT fields; adjust sizes/types if your DB requires it.
//   - CreatedAt is stored as TIMESTAMP (without time zone). Change to TIMESTAMPTZ
//     if you prefer storing timezone-aware timestamps.
//
// This SQL is written to be compatible with common SQL dialects (e.g., Postgres
// and SQLite). You may tweak data types or constraints to better match your DB.
const CreateUsersTableSQL = `
CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  created_at TIMESTAMPTZ NOT NULL
);
`

// CreateOutboxEventsTableSQL contains the SQL statement to create an outbox_events table
// that matches the OutboxEvent struct. This is written to be compatible with SQLite.
//
// Notes for SQLite:
// - Use TEXT for string-like values.
// - TIMESTAMP is accepted by SQLite and stored as text/numeric internally depending on value.
// - payload is stored as TEXT to hold JSON; metadata is optional and may be NULL.
// - processed_at is nullable to indicate events not yet processed.
const CreateOutboxEventsTableSQL = `
CREATE TABLE IF NOT EXISTS outbox_events (
  id TEXT PRIMARY KEY,
  aggregate_type TEXT NOT NULL,
  aggregate_id TEXT NOT NULL,
  event_type TEXT NOT NULL,
  payload TEXT NOT NULL,
  metadata TEXT,
  created_at TIMESTAMPTZ NOT NULL,
  processed_at TIMESTAMPTZ
);
CREATE INDEX IF NOT EXISTS idx_outbox_events_processed_at ON outbox_events (processed_at) WHERE processed_at IS NULL;
`
