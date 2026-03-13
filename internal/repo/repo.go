package repo

import (
	"context"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	systemrepo "github.com/Cheasezz/balanceSrvc/internal/repo/system"
	"github.com/Cheasezz/balanceSrvc/pkg/pgx5"
)

type System interface {
	Transaction(c context.Context, req *core.SystemTransaction) error
}

type Repo struct {
	System System
}

func New(db *pgx5.Pgx) *Repo {
	return &Repo{
		System: systemrepo.New(db),
	}
}
