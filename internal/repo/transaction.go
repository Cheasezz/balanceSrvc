package repo

import (
	"context"
	"fmt"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/Cheasezz/balanceSrvc/pkg/pgx5"
)

type trxRepo struct {
	db *pgx5.Pgx
}

func newTrxRepo(db *pgx5.Pgx) *trxRepo {
	return &trxRepo{db}
}

func (r *trxRepo) GetAllTypesInfo(ctx context.Context) (map[string]*core.TrxType, error) {
	const op = "trxrepo.GetAllTypesInfo"

	query := fmt.Sprintf("SELECT * from %s", trxTypesTable)

	rows, err := r.db.Pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	typeMap := make(map[string]*core.TrxType)

	for rows.Next() {
		var t core.TrxType
		if err := rows.Scan(&t.Id, &t.Code, &t.Name, &t.Category); err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}

		typeMap[t.Code] = &t
	}
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return typeMap, nil
}
