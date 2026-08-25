package postgres

import (
	"errors"

	"github.com/Cheasezz/balanceSrvc/pkg/pgx5"
)

type Postgres struct {
	db *pgx5.Pgx
}

const (
	userTable     = "users"
	trxTable      = "transactions"
	trxTypesTable = "transaction_types"
)

var (
	ErrInsuffBalance  = errors.New("insufficient balance")
	ErrIdNotfound     = errors.New("id not found in db")
	ErrIdempotencyKey = errors.New("idempotency key already exists")
)

func New(db *pgx5.Pgx) *Postgres {
	return &Postgres{db}
}
