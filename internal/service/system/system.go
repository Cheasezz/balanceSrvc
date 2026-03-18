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

func (s *service) TransactionTo(
	ctx context.Context,
	userId uuid.UUID,
	amount int64,
	trxType blnc.SystemTrxToType,
) error {

	const op = "systemsrvc.TransactionTo"

	tType, err := s.rg.SystemToType(trxType)
	if err != nil {
		return fmt.Errorf("op=%s, %w", op, err)
	}

	trxInfo := &core.Transaction{
		Type_id:      tType.Id,
		Resipient_id: userId,
		Amount:       amount,
	}

	s.db.System.TransactionTo(ctx, trxInfo)
	return nil
}
