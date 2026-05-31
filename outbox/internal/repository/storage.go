package repository

import "context"

// txKey is a globally accessible, type-safe key used to pass transactions
// through context boundaries across different domain implementations.
type txKey string

const (
	SqliteTxKey   txKey = "sqlite_tx"
	PostgresTxKey txKey = "postgres_tx"
)

// Transactor defines the universal API for managing transactional boundaries,
// regardless of whether the underlying storage is SQLite, Postgres, or Mongo.
type Transactor interface {
	Transact(ctx context.Context, fn func(ctx context.Context) error) error
}
