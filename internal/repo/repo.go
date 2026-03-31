package repo

import (
	"context"
	"errors"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/Cheasezz/balanceSrvc/pkg/pgx5"
	"github.com/google/uuid"
)

type System interface {
	TransactionTo(c context.Context, trx *core.Transaction) error
	TransactionFrom(c context.Context, trx *core.Transaction) error
}

type User interface {
	TransactionToUser(c context.Context, trx *core.Transaction) error
	Balance(c context.Context, userId uuid.UUID) (int64, error)
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
	ErrIdNotfound    = errors.New("id not found in db")
)

func New(db *pgx5.Pgx) *Repo {
	return &Repo{
		System: newSystemRepo(db),
		User:   newUserRepo(db),
		Trx:    newTrxRepo(db),
	}
}
