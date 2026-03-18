package grpcHndlrs

import (
	"context"
	"errors"

	trxtyperegistry "github.com/Cheasezz/balanceSrvc/internal/app/trxTypeRegistry"
	"github.com/Cheasezz/balanceSrvc/internal/service"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

const (
	errInvalidUuid    = "field user_id must be valid uuid"
	errInvalidTrxType = "unacceptable transaction type"
	errInvalidAmount  = "field amount must be int and not equal to 0"
	errInternalServer = "something went wrong on server"
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
		return nil, status.Error(codes.InvalidArgument, errInvalidUuid)
	}

	if req.GetAmount() == 0 {
		return nil, status.Error(codes.InvalidArgument, errInvalidAmount)
	}

	err = s.srvc.System.TransactionTo(ctx, id, req.GetAmount(), req.SystemTrxType)
	if err != nil {
		if errors.Is(err, trxtyperegistry.ErrUnknowSysTrxToType) {
			return nil, status.Error(codes.InvalidArgument, errInvalidTrxType)
		}
		return nil, status.Error(codes.Internal, errInternalServer)
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
