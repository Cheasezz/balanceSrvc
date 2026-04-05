package service

import (
	"context"
	"errors"

	"github.com/Cheasezz/balanceSrvc/internal/adapter/postgres"
	trxtyperegistry "github.com/Cheasezz/balanceSrvc/internal/adapter/trxTypeRegistry"
	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/Cheasezz/balanceSrvc/internal/dto"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
)

type PgSystem interface {
	TransactionTo(c context.Context, trx *core.Transaction) error
	TransactionFrom(c context.Context, trx *core.Transaction) error
}

type systemSrvc struct {
	log logger.Logger
	pg  PgSystem
	rg  trxTypeRegistry
}

func NewSystemSrvc(l logger.Logger, db PgSystem, tr trxTypeRegistry) *systemSrvc {
	return &systemSrvc{l, db, tr}
}

func (s *systemSrvc) TransactionTo(ctx context.Context, input dto.SystemTrxInput) error {

	const op = "systemsrvc.TransactionTo"
	log := s.log.With("op", op)

	tType, err := s.rg.SystemToType(input.TrxType)
	if err != nil {
		log.Error("failed to check transaction type", "err", err)

		if errors.Is(err, trxtyperegistry.ErrUnknowSysTrxToType) {
			return core.ErrUnknownTrxType
		}

		return err
	}

	trxInfo, err := core.NewSystemToUserTrx(tType, input.UserId, input.Amount)
	if err != nil {
		log.Error("failed to create new systemToUser transaction", "err", err)
		return err
	}

	err = s.pg.TransactionTo(ctx, trxInfo)
	if err != nil {
		log.Error("failed postgres method", "err", err)
		return err
	}

	return nil
}

func (s *systemSrvc) TransactionFrom(ctx context.Context, input dto.SystemTrxInput) error {

	const op = "systemsrvc.TransactionFrom"
	log := s.log.With("op", op)

	tType, err := s.rg.SystemFromType(input.TrxType)
	if err != nil {
		log.Error("failed to check transaction type", "err", err)
		if errors.Is(err, trxtyperegistry.ErrUnknowSysTrxFromType) {
			return core.ErrUnknownTrxType
		}

		return err
	}

	trxInfo, err := core.NewSystemFromUserTrx(tType, input.UserId, input.Amount)
	if err != nil {
		log.Error("failed to create new systemFromUser transaction", "err", err)
		return err
	}

	err = s.pg.TransactionFrom(ctx, trxInfo)
	if err != nil {
		log.Error("failed postgres method", "err", err)
		if errors.Is(err, postgres.ErrInsuffBalance) {
			return core.ErrInsuffBalance
		}
		return err
	}

	return nil
}
