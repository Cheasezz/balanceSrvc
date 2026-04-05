package postgres

import (
	"context"
	"fmt"

	"github.com/Cheasezz/balanceSrvc/internal/core"
)

// type systemRepo struct {
// 	db *pgx5.Pgx
// }

// func newSystemRepo(db *pgx5.Pgx) *systemRepo {
// 	return &systemRepo{db}
// }

func (r *Postgres) TransactionTo(ctx context.Context, trx *core.Transaction) error {
	const op = "systemrepo.TransactionTo"

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer tx.Rollback(ctx)

	query := fmt.Sprintf(
		`INSERT INTO %s AS u (id, balance) VALUES ($1, $2) ON CONFLICT (id) DO 
		UPDATE SET balance = u.balance + EXCLUDED.balance`,
		userTable)
	_, err = tx.Exec(ctx, query, trx.Resipient_id, trx.Amount)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	query = fmt.Sprintf(
		`INSERT INTO %s (resipient_id , type_id, amount) VALUES ($1, $2, $3)`,
		trxTable)
	_, err = tx.Exec(ctx, query, trx.Resipient_id, trx.Type_id, trx.Amount)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (r *Postgres) TransactionFrom(ctx context.Context, trx *core.Transaction) error {
	const op = "systemrepo.TransactionFrom"

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
		`INSERT INTO %s (sender_id , type_id, amount) VALUES ($1, $2, $3)`,
		trxTable)
	_, err = tx.Exec(ctx, query, trx.Sender_id, trx.Type_id, trx.Amount)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
