package core

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	Id           uuid.UUID `db:"id"`
	Sender_id    uuid.UUID `db:"sender_id"`
	Resipient_id uuid.UUID `db:"resipient_id"`
	Type_id      uint8     `db:"type_id"`
	Amount       int64     `db:"amount"`
	Created_at   time.Time `db:"created_at"`
}

type TrxType struct {
	Id       uint8  `db:"id"`
	Code     string `db:"code"`
	Name     string `db:"name"`
	Category string `db:"category"`
}
