package service

import (
	"context"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/Cheasezz/balanceSrvc/internal/repo"
	systemsrvc "github.com/Cheasezz/balanceSrvc/internal/service/system"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
)

type System interface {
	Transaction(c context.Context, req *core.SystemTransaction) error
}

type Service struct {
	System System
}

func New(l logger.Logger, db *repo.Repo) *Service {
	return &Service{
		System: systemsrvc.New(l, db),
	}
}
