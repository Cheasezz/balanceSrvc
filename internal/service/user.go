package service

import (
	"context"
	"errors"

	trxtyperegistry "github.com/Cheasezz/balanceSrvc/internal/app/trxTypeRegistry"
	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/Cheasezz/balanceSrvc/internal/repo"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
	"github.com/google/uuid"
)

var (
	ErrUsrTrxType          = errors.New("unknow user transaction type")
	ErrUserTrxTypeDisabled = errors.New("this type is disabled")
	ErrSameIds             = errors.New("Ids must be not equal")
)

type userSrvc struct {
	log logger.Logger
	db  *repo.Repo
	rg  trxTypeRegistry
}

func NewUserSrvc(l logger.Logger, db *repo.Repo, tr trxTypeRegistry) *userSrvc {
	return &userSrvc{l, db, tr}
}

func (s *userSrvc) TransactionToUser(
	ctx context.Context,
	sender,
	resipient uuid.UUID,
	amount uint64,
	trxType blnc.UserTrxType,
) error {

	const op = "usersrvc.TransactionToUser"
	log := s.log.With("op", op)

	tType, err := s.rg.UserType(trxType)
	if err != nil {
		log.Error("failed to check transaction type", "err", err)

		if errors.Is(err, trxtyperegistry.ErrUnknowUsrTrxType) {
			return ErrUsrTrxType
		}

		return err
	}

	trxInfo, err := core.NewUserToUserTrx(tType, sender, resipient, amount)
	if err != nil {
		log.Error("failed to create new UserToUser transaction", "err", err)
		switch {
		case errors.Is(err, core.ErrDisabledType):
			return ErrUserTrxTypeDisabled
		case errors.Is(err, core.ErrSameIds):
			return ErrSameIds
		}
		return err
	}

	err = s.db.User.TransactionToUser(ctx, trxInfo)
	if err != nil {
		log.Error("failed repo method", "err", err)
		if errors.Is(err, repo.ErrInsuffBalance) {
			return ErrInsuffBalance
		}
		return err
	}
	return nil
}
