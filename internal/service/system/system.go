package systemsrvc

import (
	"context"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/Cheasezz/balanceSrvc/internal/repo"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
)

type service struct {
	log logger.Logger
	db  *repo.Repo
}

func New(l logger.Logger, db *repo.Repo) *service {
	return &service{l, db}
}

func (s *service) Transaction(ctx context.Context, req *core.SystemTransaction) error {
	s.db.System.Transaction(ctx, req)
	return nil
}
