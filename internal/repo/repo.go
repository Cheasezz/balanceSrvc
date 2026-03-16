package repo

import (
	"context"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	systemrepo "github.com/Cheasezz/balanceSrvc/internal/repo/system"
	trxrepo "github.com/Cheasezz/balanceSrvc/internal/repo/transaction"
	"github.com/Cheasezz/balanceSrvc/pkg/pgx5"
)

type System interface {
	Transaction(c context.Context, req *core.Transaction) error
}

type Transaction interface {
	GetAllTypesInfo(ctx context.Context) (map[string]*core.TrxType, error)
}

type Repo struct {
	System System
	Trx    Transaction
}

func New(db *pgx5.Pgx) *Repo {
	return &Repo{
		System: systemrepo.New(db),
		Trx:    trxrepo.New(db),
	}
}
