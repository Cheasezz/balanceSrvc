package systemrepo

import (
	"context"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/Cheasezz/balanceSrvc/pkg/pgx5"
)

type repo struct {
	db *pgx5.Pgx
}

func New(db *pgx5.Pgx) *repo {
	return &repo{db}
}

func (r *repo) Transaction(ctx context.Context, req *core.SystemTransaction) error {
	return nil
}
