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

func NewSystemToUserTrx(trxType *TrxType, userId string, amount uint64) (*Transaction, error) {
	id, err := sysValid(trxType, userId, amount)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		Resipient_id: id,
		Type_id:      trxType.Id,
		Amount:       amount,
	}, nil
}

func NewSystemFromUserTrx(trxType *TrxType, userId string, amount uint64) (*Transaction, error) {
	id, err := sysValid(trxType, userId, amount)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		Sender_id: id,
		Type_id:   trxType.Id,
		Amount:    amount,
	}, nil
}

func NewUserToUserTrx(trxType *TrxType, sender, resipient string, amount uint64) (*Transaction, error) {
	send, resip, err := usrValid(trxType, sender, resipient, amount)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		Sender_id:    send,
		Resipient_id: resip,
		Type_id:      trxType.Id,
		Amount:       amount,
	}, nil
}

func sysValid(trxType *TrxType, userId string, amount uint64) (uuid.UUID, error) {
	emptyVal := uuid.UUID{}
	if !trxType.Enable {
		return emptyVal, ErrDisabledType
	}

	if trxType.Category != "system" {
		return emptyVal, ErrInvalidTrxCategory
	}

	id, err := uuid.Parse(userId)
	if err != nil {
		return emptyVal, ErrInvalidUuid
	}

	if amount <= 0 {
		return emptyVal, ErrInvalidAmount
	}
	return id, nil
}

func usrValid(trxType *TrxType, sender, resipient string, amount uint64) (uuid.UUID, uuid.UUID, error) {
	emptyVal := uuid.UUID{}

	if !trxType.Enable {
		return emptyVal, emptyVal, ErrDisabledType
	}

	if trxType.Category != "user" {
		return emptyVal, emptyVal, ErrInvalidTrxCategory
	}

	send, err := uuid.Parse(sender)
	if err != nil {
		return emptyVal, emptyVal, ErrInvalidUuid
	}

	resip, err := uuid.Parse(resipient)
	if err != nil {
		return emptyVal, emptyVal, ErrInvalidUuid
	}

	if sender == resipient {
		return emptyVal, emptyVal, ErrSameIds
	}

	if amount <= 0 {
		return emptyVal, emptyVal, ErrInvalidAmount
	}
	return send, resip, nil
}
