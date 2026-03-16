package systemsrvc

import (
	"context"
	"fmt"

	trxtyperegistry "github.com/Cheasezz/balanceSrvc/internal/app/trxTypeRegistry"
	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/Cheasezz/balanceSrvc/internal/repo"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
	"github.com/google/uuid"
)

type service struct {
	log logger.Logger
	db  *repo.Repo
	rg  *trxtyperegistry.Registry
}

func New(l logger.Logger, db *repo.Repo, tr *trxtyperegistry.Registry) *service {
	return &service{l, db, tr}
}

func (s *service) Transaction(
	ctx context.Context,
	userId uuid.UUID,
	amount int64,
	trxType blnc.SystemTrxType,
) error {

	const op = "systemsrvc.Transaction"

	tType, err := s.rg.SystemType(trxType)
	if err != nil {
		return fmt.Errorf("op=%s, err=", op, err)
	}

	trxInfo := &core.Transaction{
		Type_id: tType.Id,
		Amount:  amount,
	}
	if amount < 0 {
		trxInfo.Sender_id = userId
	} else {
		trxInfo.Resipient_id = userId
	}

	s.db.System.Transaction(ctx, trxInfo)
	return nil
}
