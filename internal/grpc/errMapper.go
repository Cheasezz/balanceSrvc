package grpcSrv

import (
	"errors"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TODO: Мапить ошибки и код, а не статус целиком.
func toStatus(err error) error {
	switch {
	case errors.Is(err, core.ErrUnknownTrxType):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, core.ErrDisabledType):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, core.ErrInvalidAmount):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, core.ErrInsuffBalance):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, core.ErrInvalidUuid):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, core.ErrSameIds):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, core.ErrIdNotfound):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, core.ErrInternalServer.Error())
	}
}
