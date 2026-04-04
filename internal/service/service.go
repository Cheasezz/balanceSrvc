package service

import (
	"context"
	"errors"

	"github.com/Cheasezz/balanceSrvc/internal/adapter/postgres"
	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
	"github.com/google/uuid"
)

var (
	ErrInsuffBalance = errors.New("insufficient balance")
)

type System interface {
	TransactionTo(
		ctx context.Context,
		userId uuid.UUID,
		amount uint64,
		trxType blnc.SystemTrxToType,
	) error

	TransactionFrom(
		ctx context.Context,
		userId uuid.UUID,
		amount uint64,
		trxType blnc.SystemTrxFromType,
	) error
}

type User interface {
	TransactionToUser(
		ctx context.Context,
		sender,
		resipient uuid.UUID,
		amount uint64,
		trxType blnc.UserTrxType,
	) error

	Balance(c context.Context, userId uuid.UUID) (int64, error)
}

type trxTypeRegistry interface {
	SystemToType(t blnc.SystemTrxToType) (*core.TrxType, error)
	SystemFromType(t blnc.SystemTrxFromType) (*core.TrxType, error)
	UserType(t blnc.UserTrxType) (*core.TrxType, error)
}

type Service struct {
	System System
	User   User
}

func New(l logger.Logger, db *postgres.Postgres, tr trxTypeRegistry) *Service {
	return &Service{
		System: NewSystemSrvc(l, db, tr),
		User:   NewUserSrvc(l, db, tr),
	}
}
