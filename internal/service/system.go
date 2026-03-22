package service

import (
	"context"
	"errors"
	"fmt"

	trxtyperegistry "github.com/Cheasezz/balanceSrvc/internal/app/trxTypeRegistry"
	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/Cheasezz/balanceSrvc/internal/repo"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
	"github.com/google/uuid"
)

var (
	ErrSystemTrxToType = errors.New("unknow system transaction(to) type")
)

type systemSrvc struct {
	log logger.Logger
	db  *repo.Repo
	rg  trxTypeRegistry
}

func NewSystemSrvc(l logger.Logger, db *repo.Repo, tr trxTypeRegistry) *systemSrvc {
	return &systemSrvc{l, db, tr}
}

func (s *systemSrvc) TransactionTo(
	ctx context.Context,
	userId uuid.UUID,
	amount uint64,
	trxType blnc.SystemTrxToType,
) error {

	const op = "systemsrvc.TransactionTo"
	log := s.log.With("op", op)

	tType, err := s.rg.SystemToType(trxType)
	if err != nil {
		log.Error("failed to check transaction type", "err", err)

		if errors.Is(err, trxtyperegistry.ErrUnknowSysTrxToType) {
			return fmt.Errorf("%s: %w", op, ErrSystemTrxToType)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	trxInfo, err := core.NewSystemToUserTrx(tType, userId, amount)
	if err != nil {
		log.Error("failed to create new systemToUser transaction", "err", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	err = s.db.System.TransactionTo(ctx, trxInfo)
	if err != nil {
		log.Error("failed repo method", "err", err)
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
