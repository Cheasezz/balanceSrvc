package grpcHndlrs

import (
	"context"
	"errors"

	"github.com/Cheasezz/balanceSrvc/internal/service"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *ServerAPI) UserTransaction(
	ctx context.Context,
	req *blnc.UserTrxRequest,
) (*blnc.UserTrxResponse, error) {
	sender, err := uuid.Parse(req.GetSenderId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, ErrInvalidUuid.Error())
	}
	resipient, err := uuid.Parse(req.GetResipientId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, ErrInvalidUuid.Error())
	}

	if req.GetAmount() == 0 {
		return nil, status.Error(codes.InvalidArgument, ErrInvalidAmount.Error())
	}

	err = s.Srvc.User.TransactionToUser(ctx, sender, resipient, req.GetAmount(), req.UserTrxType)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrUsrTrxType):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, service.ErrUserTrxTypeDisabled):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, service.ErrInsuffBalance):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, service.ErrSameIds):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, ErrInternalServer.Error())
		}
	}

	return &blnc.UserTrxResponse{}, nil
}
