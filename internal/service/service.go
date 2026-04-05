package service

import (
	"context"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/Cheasezz/balanceSrvc/internal/dto"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
)

type pgDb interface {
	PgSystem
	PgUser
}

type trxTypeRegistry interface {
	SystemToType(t int32) (*core.TrxType, error)
	SystemFromType(t int32) (*core.TrxType, error)
	UserType(t int32) (*core.TrxType, error)
}

type System interface {
	TransactionTo(ctx context.Context, input dto.SystemTrxInput) error
	TransactionFrom(ctx context.Context, input dto.SystemTrxInput) error
}

type User interface {
	TransactionToUser(ctx context.Context, input dto.UserTrxInput) error
	Balance(c context.Context, userId string) (uint64, error)
}

type Service struct {
	System System
	User   User
}

func New(l logger.Logger, db pgDb, tr trxTypeRegistry) *Service {
	return &Service{
		System: NewSystemSrvc(l, db, tr),
		User:   NewUserSrvc(l, db, tr),
	}
}
