package service

import (
	"context"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/Cheasezz/balanceSrvc/internal/repo"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
	"github.com/google/uuid"
)

type System interface {
	TransactionTo(
		ctx context.Context,
		userId uuid.UUID,
		amount uint64,
		trxType blnc.SystemTrxToType,
	) error
}

type trxTypeRegistry interface {
	SystemToType(t blnc.SystemTrxToType) (*core.TrxType, error)
	SystemFromType(t blnc.SystemTrxFromType) (*core.TrxType, error)
	UserType(t blnc.UserTrxType) (*core.TrxType, error)
}

type Service struct {
	System *systemSrvc
}

func New(l logger.Logger, db *repo.Repo, tr trxTypeRegistry) *Service {
	return &Service{
		System: NewSystemSrvc(l, db, tr),
	}
}
