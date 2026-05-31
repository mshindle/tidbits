package sqlite

import (
	"context"
	"database/sql"
	"errors"

	"github.com/apex/log"
	"gitlab.com/mshindle/tidbits/outbox/internal/repository"
)

func (s *Store) Transact(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer handleRollback(tx)

	// Inject the standard *sql.Tx using the global contract key
	txCtx := context.WithValue(ctx, repository.SqliteTxKey, tx)

	if err := fn(txCtx); err != nil {
		return err
	}
	return tx.Commit()
}

func handleRollback(tx *sql.Tx) {
	if rbe := tx.Rollback(); rbe != nil && !errors.Is(rbe, sql.ErrTxDone) {
		log.WithError(rbe).Fatal("rolling back transaction failed")
	}
}
