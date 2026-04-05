package grpcSrv

import (
	"context"

	"github.com/Cheasezz/balanceSrvc/internal/dto"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
)

func (s *ServerAPI) UserTransaction(
	ctx context.Context,
	req *blnc.UserTrxRequest,
) (*blnc.UserTrxResponse, error) {

	input := dto.UserTrxInput{
		Sender:    req.GetSenderId(),
		Resipient: req.GetResipientId(),
		Amount:    req.GetAmount(),
		TrxType:   int32(req.UserTrxType),
	}
	err := s.Srvc.User.TransactionToUser(ctx, input)
	if err != nil {
		return nil, toStatus(err)
	}

	return &blnc.UserTrxResponse{}, nil
}

func (s *ServerAPI) UserBalance(
	ctx context.Context,
	req *blnc.BalanceRequest,
) (*blnc.BalanceResponse, error) {

	balance, err := s.Srvc.User.Balance(ctx, req.GetUserId())
	if err != nil {
		return nil, toStatus(err)
	}

	return &blnc.BalanceResponse{Balance: balance}, nil
}
