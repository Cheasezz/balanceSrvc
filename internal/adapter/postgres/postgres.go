package postgres

import (
	"errors"

	"github.com/Cheasezz/balanceSrvc/pkg/pgx5"
)

// type System interface {
// 	TransactionTo(c context.Context, trx *core.Transaction) error
// 	TransactionFrom(c context.Context, trx *core.Transaction) error
// }

// type User interface {
// 	TransactionToUser(c context.Context, trx *core.Transaction) error
// 	Balance(c context.Context, userId uuid.UUID) (uint64, error)
// }

// type Transaction interface {
// 	GetAllTypesInfo(ctx context.Context) (map[string]*core.TrxType, error)
// }

type Postgres struct {
	db *pgx5.Pgx
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

func New(db *pgx5.Pgx) *Postgres {
	return &Postgres{db}
}
