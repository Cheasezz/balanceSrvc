package grpcHndlrs

import (
	"errors"

	"github.com/Cheasezz/balanceSrvc/internal/service"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	ErrInvalidUuid = errors.New("field user_id must be valid uuid")
	// ErrInvalidTrxType  = errors.New("unacceptable transaction type")
	// ErrTrxTypeDisabled = errors.New("transaction with this type doesent accept at this moment")
	ErrInvalidAmount  = errors.New("field amount must be uint64 and not equal to 0")
	ErrInternalServer = errors.New("something went wrong on server")
	ErrIdNotFound     = errors.New("id not found")
	// ErrInsuffBalance   = errors.New("insufficient balance")
)

const (
	envLocal = "local"
)

type ServerAPI struct {
	blnc.UnimplementedBalanceServer
	Srvc *service.Service
}

func Register(gRPC *grpc.Server, l logger.Logger, s *service.Service, env string) {
	blnc.RegisterBalanceServer(gRPC, &ServerAPI{Srvc: s})
	if env == envLocal {
		reflection.Register(gRPC)
	}
}
