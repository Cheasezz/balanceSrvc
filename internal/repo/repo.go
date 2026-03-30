package repo

import (
	"context"
	"errors"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/Cheasezz/balanceSrvc/pkg/pgx5"
)

type System interface {
	TransactionTo(c context.Context, trx *core.Transaction) error
	TransactionFrom(c context.Context, trx *core.Transaction) error
}

type User interface {
	TransactionToUser(c context.Context, trx *core.Transaction) error
}

type Transaction interface {
	GetAllTypesInfo(ctx context.Context) (map[string]*core.TrxType, error)
}

type Repo struct {
	System System
	User   User
	Trx    Transaction
}

const (
	userTable     = "users"
	trxTable      = "transactions"
	trxTypesTable = "transaction_types"
)

var (
	ErrInsuffBalance = errors.New("insufficient balance")
)

func New(db *pgx5.Pgx) *Repo {
	return &Repo{
		System: newSystemRepo(db),
		User:   newUserRepo(db),
		Trx:    newTrxRepo(db),
	}
}
