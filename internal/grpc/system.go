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

func (s *ServerAPI) SystemTransactionTo(
	ctx context.Context,
	req *blnc.SystemTrxToRequest,
) (*blnc.SystemTrxResponse, error) {

	id, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, ErrInvalidUuid.Error())
	}

	if req.GetAmount() == 0 {
		return nil, status.Error(codes.InvalidArgument, ErrInvalidAmount.Error())
	}

	err = s.Srvc.System.TransactionTo(ctx, id, req.GetAmount(), req.SystemTrxType)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrSystemTrxToType):
			return nil, status.Error(codes.InvalidArgument, ErrInvalidTrxType.Error())
		case errors.Is(err, service.ErrSystemTrxTypeDisabled):
			return nil, status.Error(codes.InvalidArgument, ErrSystemTrxTypeDisabled.Error())
		default:
			return nil, status.Error(codes.Internal, ErrInternalServer.Error())
		}
	}

	return &blnc.SystemTrxResponse{}, nil
}
