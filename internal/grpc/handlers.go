package grpcHndlrs

import (
	"context"
	"errors"

	"github.com/Cheasezz/balanceSrvc/internal/service"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	ErrInvalidUuid           = errors.New("field user_id must be valid uuid")
	ErrInvalidTrxType        = errors.New("unacceptable transaction type")
	ErrSystemTrxTypeDisabled = errors.New("transaction with this type doesent accept at this moment")
	ErrInvalidAmount         = errors.New("field amount must be uint64 and not equal to 0")
	ErrInternalServer        = errors.New("something went wrong on server")
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

func (s *ServerAPI) UserTransaction(
	ctx context.Context,
	req *blnc.UserTrxRequest,
) (*blnc.UserTrxResponse, error) {
	panic("Implement me pls")
	//...
}

func (s *ServerAPI) UserBalance(
	ctx context.Context,
	req *blnc.BalanceRequest,
) (*blnc.BalanceResponse, error) {
	panic("Implement me pls")
	//...
}
