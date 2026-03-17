package service

import (
	"context"

	trxtyperegistry "github.com/Cheasezz/balanceSrvc/internal/app/trxTypeRegistry"
	"github.com/Cheasezz/balanceSrvc/internal/repo"
	systemsrvc "github.com/Cheasezz/balanceSrvc/internal/service/system"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
	"github.com/google/uuid"
)

type System interface {
	TransactionTo(
		ctx context.Context,
		userId uuid.UUID,
		amount int64,
		trxType blnc.SystemTrxToType,
	) error
}

type Service struct {
	System System
}

func New(l logger.Logger, db *repo.Repo, tr *trxtyperegistry.Registry) *Service {
	return &Service{
		System: systemsrvc.New(l, db, tr),
	}
}
