package dual

import (
	"context"
	"database/sql"
	"errors"

	"github.com/apex/log"
	"gitlab.com/mshindle/tidbits/outbox/internal/repository"
)

type Transactor struct {
	sqliteDB   *sql.DB
	postgresDB *sql.DB
}

func (dt *Transactor) WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error {
	// 1. Start SQLite Transaction
	sqTx, err := dt.sqliteDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer handleRollback(sqTx) // Safe to call even if committed

	// 2. Start Postgres Transaction
	pgTx, err := dt.postgresDB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer handleRollback(pgTx)

	// 3. Inject BOTH transactions into the context
	txCtx := context.WithValue(ctx, repository.SqliteTxKey, sqTx)
	txCtx = context.WithValue(txCtx, repository.PostgresTxKey, pgTx)

	// 4. Run the domain logic
	if err := fn(txCtx); err != nil {
		return err // Both will roll back due to defers
	}

	// 5. Commit both (Best-Effort 1PC)
	if err := sqTx.Commit(); err != nil {
		return err
	}
	if err := pgTx.Commit(); err != nil {
		// WARNING: If SQLite commits but Postgres fails here,
		// you have a partial failure. This is why a shared database is preferred for the Outbox Pattern!
		return err
	}

	return nil
}

func handleRollback(tx *sql.Tx) {
	if rbe := tx.Rollback(); rbe != nil && !errors.Is(rbe, sql.ErrTxDone) {
		log.WithError(rbe).Fatal("rolling back transaction failed")
	}
}
