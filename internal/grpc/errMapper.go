package grpcSrv

import (
	"errors"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toStatus(err error) error {
	var code codes.Code

	switch {
	case errors.Is(err, core.ErrUnknownTrxType):
		code = codes.InvalidArgument
	case errors.Is(err, core.ErrDisabledType):
		code = codes.InvalidArgument
	case errors.Is(err, core.ErrInvalidAmount):
		code = codes.InvalidArgument
	case errors.Is(err, core.ErrInsuffBalance):
		code = codes.InvalidArgument
	case errors.Is(err, core.ErrInvalidUuid):
		code = codes.InvalidArgument
	case errors.Is(err, core.ErrSameIds):
		code = codes.InvalidArgument
	case errors.Is(err, core.ErrIdNotfound):
		code = codes.NotFound
	case errors.Is(err, core.ErrInvalidIdempotencyKey):
		code = codes.InvalidArgument
	case errors.Is(err, core.ErrRequestInProcess):
		code = codes.AlreadyExists
	default:
		return status.Error(codes.Internal, core.ErrInternalServer.Error())
	}

	return status.Error(code, err.Error())
}
