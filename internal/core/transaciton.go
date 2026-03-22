package core

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrDisabledType       = errors.New("this type is disabled")
	ErrInvalidAmount      = errors.New("invalid amount value")
	ErrInvalidTrxCategory = errors.New("invalid transaction category")
	ErrInvalidUserId      = errors.New("invalid user id")
)

type Transaction struct {
	Id           uuid.UUID `db:"id"`
	Sender_id    uuid.UUID `db:"sender_id"`
	Resipient_id uuid.UUID `db:"resipient_id"`
	Type_id      uint8     `db:"type_id"`
	Amount       uint64    `db:"amount"`
	Created_at   time.Time `db:"created_at"`
}

type TrxType struct {
	Id       uint8  `db:"id"`
	Code     string `db:"code"`
	Name     string `db:"name"`
	Category string `db:"category"`
	Enable   bool   `db:"enable"`
}
