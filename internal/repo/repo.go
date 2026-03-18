package repo

import (
	"context"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/Cheasezz/balanceSrvc/pkg/pgx5"
)

type System interface {
	TransactionTo(c context.Context, trx *core.Transaction) error
}

type Transaction interface {
	GetAllTypesInfo(ctx context.Context) (map[string]*core.TrxType, error)
}

type Repo struct {
	System System
	Trx    Transaction
}

const (
	userTable     = "users"
	trxTable      = "transactions"
	trxTypesTable = "transaction_types"
)

func New(db *pgx5.Pgx) *Repo {
	return &Repo{
		System: newSystemRepo(db),
		Trx:    newTrxRepo(db),
	}
}
