package grpcHndlrs

import (
	"context"
	"errors"

	"github.com/Cheasezz/balanceSrvc/internal/service"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

var (
	ErrInvalidUuid    = errors.New("field user_id must be valid uuid")
	ErrInvalidTrxType = errors.New("unacceptable transaction type")
	ErrInvalidAmount  = errors.New("field amount must be int and not equal to 0")
	ErrInternalServer = errors.New("something went wrong on server")
)

const (
	envLocal = "local"
)

type serverAPI struct {
	blnc.UnimplementedBalanceServer
	srvc *service.Service
}

func Register(gRPC *grpc.Server, l logger.Logger, s *service.Service, env string) {
	blnc.RegisterBalanceServer(gRPC, &serverAPI{srvc: s})
	if env == envLocal {
		reflection.Register(gRPC)
	}
}

func (s *serverAPI) SystemTransactionTo(
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

	err = s.srvc.System.TransactionTo(ctx, id, req.GetAmount(), req.SystemTrxType)
	if err != nil {
		if errors.Is(err, service.ErrSystemTrxToType) {
			return nil, status.Error(codes.InvalidArgument, ErrInvalidTrxType.Error())
		}
		return nil, status.Error(codes.Internal, ErrInternalServer.Error())
	}

	return &blnc.SystemTrxResponse{}, nil
}

func (s *serverAPI) UserTransaction(
	ctx context.Context,
	req *blnc.UserTrxRequest,
) (*blnc.UserTrxResponse, error) {
	panic("Implement me pls")
	//...
}

func (s *serverAPI) UserBalance(
	ctx context.Context,
	req *blnc.BalanceRequest,
) (*blnc.BalanceResponse, error) {
	panic("Implement me pls")
	//...
}
