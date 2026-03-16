package trxrepo

import (
	"context"
	"fmt"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/Cheasezz/balanceSrvc/pkg/pgx5"
)

type repo struct {
	db *pgx5.Pgx
}

const trxTable = "transaction_types"

func New(db *pgx5.Pgx) *repo {
	return &repo{db}
}

func (r *repo) GetAllTypesInfo(ctx context.Context) (map[string]*core.TrxType, error) {
	const op = "trxrepo.GetAllTypesInfo"

	query := fmt.Sprintf("SELECT * from %s", trxTable)

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("op=%s, err=%w", op, err)
	}
	defer rows.Close()

	typeMap := make(map[string]*core.TrxType)

	for rows.Next() {
		var t core.TrxType
		if err := rows.Scan(&t.Id, &t.Code, &t.Name, &t.Category); err != nil {
			return nil, fmt.Errorf("op=%s, err=%w", op, err)
		}

		typeMap[t.Code] = &t
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("op=%s, err=%w", op, err)
	}

	return typeMap, nil
}
