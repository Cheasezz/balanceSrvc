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

func (s *ServerAPI) SystemTransactionTo(
	ctx context.Context,
	req *blnc.SystemTrxToRequest,
) (*blnc.SystemTrxResponse, error) {

	input := dto.SystemTrxInput{
		UserId:  req.GetUserId(),
		Amount:  req.GetAmount(),
		TrxType: int32(req.SystemTrxType),
	}

	err := s.Srvc.System.TransactionTo(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrUnknownTrxType):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, core.ErrDisabledType):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, core.ErrInvalidAmount):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, core.ErrInvalidUuid):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, core.ErrInternalServer.Error())
		}
	}

	return &blnc.SystemTrxResponse{}, nil
}

func (s *ServerAPI) SystemTransactionFrom(
	ctx context.Context,
	req *blnc.SystemTrxFromRequest,
) (*blnc.SystemTrxResponse, error) {

	input := dto.SystemTrxInput{
		UserId:  req.GetUserId(),
		Amount:  req.GetAmount(),
		TrxType: int32(req.SystemTrxType),
	}

	err := s.Srvc.System.TransactionFrom(ctx, input)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrUnknownTrxType):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, core.ErrDisabledType):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, core.ErrInvalidAmount):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, core.ErrInsuffBalance):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		case errors.Is(err, core.ErrInvalidUuid):
			return nil, status.Error(codes.InvalidArgument, err.Error())
		default:
			return nil, status.Error(codes.Internal, core.ErrInternalServer.Error())
		}
	}

	return &blnc.SystemTrxResponse{}, nil
}
