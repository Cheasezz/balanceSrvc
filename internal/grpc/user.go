package grpcSrv

import (
	"context"
	"errors"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/Cheasezz/balanceSrvc/internal/dto"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		switch {
		case errors.Is(err, core.ErrUnknownTrxType):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, core.ErrDisabledType):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, core.ErrInsuffBalance):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, core.ErrSameIds):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, core.ErrInvalidUuid):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, core.ErrInvalidAmount):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, core.ErrInternalServer.Error())
		}
	}

	return &blnc.UserTrxResponse{}, nil
}

func (s *ServerAPI) UserBalance(
	ctx context.Context,
	req *blnc.BalanceRequest,
) (*blnc.BalanceResponse, error) {

	balance, err := s.Srvc.User.Balance(ctx, req.GetUserId())
	if err != nil {
		switch {
		case errors.Is(err, core.ErrIdNotfound):
			return nil, status.Error(codes.NotFound, err.Error())
		case errors.Is(err, core.ErrInvalidUuid):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, core.ErrInternalServer.Error())
		}
	}

	return &blnc.BalanceResponse{Balance: balance}, nil
}
