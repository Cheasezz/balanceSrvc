package repo

import (
	"context"
	"fmt"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/Cheasezz/balanceSrvc/pkg/pgx5"
)

type systemRepo struct {
	db *pgx5.Pgx
}

func newSystemRepo(db *pgx5.Pgx) *systemRepo {
	return &systemRepo{db}
}

func (r *systemRepo) TransactionTo(ctx context.Context, trx *core.Transaction) error {
	const op = "systemrepo.TransactionTo"

	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("op=%s, err=%w", op, err)
	}
	defer tx.Rollback(ctx)

	query := fmt.Sprintf(
		`INSERT INTO %s AS u (id, balance) VALUES ($1, $2) ON CONFLICT (id) DO 
		UPDATE SET balance = u.balance + EXCLUDED.balance`,
		userTable)
	_, err = tx.Exec(ctx, query, trx.Resipient_id, trx.Amount)
	if err != nil {
		return fmt.Errorf("op=%s, err=%w", op, err)
	}

	query = fmt.Sprintf(
		`INSERT INTO %s (resipient_id , type_id, amount) VALUES ($1, $2, $3)`,
		trxTable)
	_, err = tx.Exec(ctx, query, trx.Resipient_id, trx.Type_id, trx.Amount)
	if err != nil {
		return fmt.Errorf("op=%s, err=%w", op, err)
	}

	if err = tx.Commit(ctx); err != nil {
		return fmt.Errorf("op=%s, err=%w", op, err)
	}

	return nil
}
