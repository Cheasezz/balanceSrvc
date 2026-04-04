package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/Cheasezz/balanceSrvc/pkg/pgx5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type userRepo struct {
	db *pgx5.Pgx
}

func newUserRepo(db *pgx5.Pgx) *userRepo {
	return &userRepo{db}
}

func (r *userRepo) TransactionToUser(ctx context.Context, trx *core.Transaction) error {
	const op = "userrepo.TransactionToUser"

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(ctx)

	query := fmt.Sprintf(
		`UPDATE %s SET balance = balance - $1 WHERE id = $2 AND balance >= $1`,
		userTable)
	ct, err := tx.Exec(ctx, query, trx.Amount, trx.Sender_id)

	if ct.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", op, ErrInsuffBalance)
	}
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query = fmt.Sprintf(
		`INSERT INTO %s AS u (id, balance) VALUES ($1, $2) ON CONFLICT (id) DO 
		UPDATE SET balance = u.balance + EXCLUDED.balance`,
		userTable)
	_, err = tx.Exec(ctx, query, trx.Resipient_id, trx.Amount)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query = fmt.Sprintf(
		`INSERT INTO %s (sender_id, resipient_id, type_id, amount) VALUES ($1, $2, $3, $4)`,
		trxTable)
	_, err = tx.Exec(ctx, query, trx.Sender_id, trx.Resipient_id, trx.Type_id, trx.Amount)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}

func (r *userRepo) Balance(ctx context.Context, userId uuid.UUID) (int64, error) {
	const op = "userrepo.Balance"
	var balance int64

	query := fmt.Sprintf("SELECT balance from %s WHERE id = $1", userTable)

	err := r.db.Scany.Get(ctx, r.db.Pool, &balance, query, userId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return balance, fmt.Errorf("%s: %w", op, ErrIdNotfound)
		}
		return balance, err
	}

	return balance, nil
}
