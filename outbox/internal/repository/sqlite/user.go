package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"gitlab.com/mshindle/tidbits/outbox/internal/repository"
	"gitlab.com/mshindle/tidbits/outbox/internal/user"
)

func (s *Store) CreateUser(ctx context.Context, u *user.User) error {
	stmt := `INSERT INTO users (name, email, created_at) VALUES ($1, $2, $3)`

	var err error
	var res sql.Result

	// Intercept the global transaction key provided by the contract package
	if tx, ok := ctx.Value(repository.SqliteTxKey).(*sql.Tx); ok {
		res, err = tx.ExecContext(ctx, stmt, u.Name, u.Email, u.CreatedAt)
	} else {
		res, err = s.db.ExecContext(ctx, stmt, u.Name, u.Email, u.CreatedAt)
	}

	if err != nil {
		return fmt.Errorf("failed to insert user: %w", err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("failed to retrieve last insert id: %w", err)
	}

	u.ID = id
	return nil
}
