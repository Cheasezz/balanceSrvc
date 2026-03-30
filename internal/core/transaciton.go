package core

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrDisabledType       = errors.New("this type is disabled")
	ErrInvalidAmount      = errors.New("invalid amount value, must be uint and not equal to 0")
	ErrInvalidTrxCategory = errors.New("this type of transaction not in current catagory")
	ErrInvalidUserId      = errors.New("invalid user id (uuid.Nil)")
	ErrSameIds            = errors.New("Ids must be not equal")
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

func NewSystemToUserTrx(trxType *TrxType, userId uuid.UUID, amount uint64) (*Transaction, error) {
	err := sysValid(trxType, userId, amount)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		Resipient_id: userId,
		Type_id:      trxType.Id,
		Amount:       amount,
	}, nil
}

func NewSystemFromUserTrx(trxType *TrxType, userId uuid.UUID, amount uint64) (*Transaction, error) {
	err := sysValid(trxType, userId, amount)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		Sender_id: userId,
		Type_id:   trxType.Id,
		Amount:    amount,
	}, nil
}

func NewUserToUserTrx(trxType *TrxType, sender, resipient uuid.UUID, amount uint64) (*Transaction, error) {
	err := usrValid(trxType, sender, resipient, amount)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		Sender_id:    sender,
		Resipient_id: resipient,
		Type_id:      trxType.Id,
		Amount:       amount,
	}, nil
}

func sysValid(trxType *TrxType, userId uuid.UUID, amount uint64) error {
	if !trxType.Enable {
		return ErrDisabledType
	}

	if trxType.Category != "system" {
		return ErrInvalidTrxCategory
	}

	if userId == uuid.Nil {
		return ErrInvalidUserId
	}

	if amount <= 0 {
		return ErrInvalidAmount
	}
	return nil
}

func usrValid(trxType *TrxType, sender, resipient uuid.UUID, amount uint64) error {
	if !trxType.Enable {
		return ErrDisabledType
	}

	if trxType.Category != "user" {
		return ErrInvalidTrxCategory
	}

	if sender == uuid.Nil || resipient == uuid.Nil {
		return ErrInvalidUserId
	}

	if sender == resipient {
		return ErrSameIds
	}

	if amount <= 0 {
		return ErrInvalidAmount
	}
	return nil
}
