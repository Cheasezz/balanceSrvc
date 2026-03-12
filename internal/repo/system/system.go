package systemrepo

import (
	"context"

	"github.com/Cheasezz/balanceSrvc/internal/core"
)

type repo struct {
}

func New() *repo {
	return &repo{}
}

func (r *repo) Transaction(ctx context.Context, req *core.SystemTransaction) error {
	return nil
}
