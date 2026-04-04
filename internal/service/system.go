package service

import (
	"context"
	"errors"

	"github.com/Cheasezz/balanceSrvc/internal/adapter/postgres"
	trxtyperegistry "github.com/Cheasezz/balanceSrvc/internal/app/trxTypeRegistry"
	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
	"github.com/google/uuid"
)

var (
	ErrSystemTrxToType       = errors.New("unknow system transaction(to) type")
	ErrSystemTrxFromType     = errors.New("unknow system transaction(from) type")
	ErrSystemTrxTypeDisabled = errors.New("this type is disabled")
)

type systemSrvc struct {
	log logger.Logger
	pg  *postgres.Postgres
	rg  trxTypeRegistry
}

func NewSystemSrvc(l logger.Logger, db *postgres.Postgres, tr trxTypeRegistry) *systemSrvc {
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
			return ErrSystemTrxToType
		}

		return err
	}

	trxInfo, err := core.NewSystemToUserTrx(tType, userId, amount)
	if err != nil {
		log.Error("failed to create new systemToUser transaction", "err", err)
		if errors.Is(err, core.ErrDisabledType) {
			return ErrSystemTrxTypeDisabled
		}
		return err
	}

	err = s.pg.System.TransactionTo(ctx, trxInfo)
	if err != nil {
		log.Error("failed postgres method", "err", err)
		return err
	}

	return nil
}

func (s *systemSrvc) TransactionFrom(
	ctx context.Context,
	userId uuid.UUID,
	amount uint64,
	trxType blnc.SystemTrxFromType,
) error {

	const op = "systemsrvc.TransactionFrom"
	log := s.log.With("op", op)

	tType, err := s.rg.SystemFromType(trxType)
	if err != nil {
		log.Error("failed to check transaction type", "err", err)

		if errors.Is(err, trxtyperegistry.ErrUnknowSysTrxFromType) {
			return ErrSystemTrxFromType
		}

		return err
	}

	trxInfo, err := core.NewSystemFromUserTrx(tType, userId, amount)
	if err != nil {
		log.Error("failed to create new systemFromUser transaction", "err", err)
		if errors.Is(err, core.ErrDisabledType) {
			return ErrSystemTrxTypeDisabled
		}
		return err
	}

	err = s.pg.System.TransactionFrom(ctx, trxInfo)
	if err != nil {
		log.Error("failed postgres method", "err", err)
		if errors.Is(err, postgres.ErrInsuffBalance) {
			return ErrInsuffBalance
		}
		return err
	}

	return nil
}
