package core

import "github.com/google/uuid"

type SystemTransaction struct {
	UserId          uuid.UUID
	TransactionType string
	Amount          int64
}
