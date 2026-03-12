package grpcHndlrs

import (
	"context"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/Cheasezz/balanceSrvc/internal/service"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	errInvalidUuid    = "field user_id must be valid uuid"
	errInvalidAmount  = "field amount must be int and not equal to 0"
	errInternalServer = "something went wrong on server"
)

type serverAPI struct {
	blnc.UnimplementedBalanceServer
	srvc *service.Service
}

func Register(gRPC *grpc.Server, s *service.Service) {
	blnc.RegisterBalanceServer(gRPC, &serverAPI{srvc: s})
}

func (s *serverAPI) SystemTransaction(
	ctx context.Context,
	req *blnc.SystemTransactionRequest,
) (*blnc.SystemTransactionResponse, error) {

	id, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, errInvalidUuid)
	}

	// slices.Contains()

	if req.GetAmount() == 0 {
		return nil, status.Error(codes.InvalidArgument, errInvalidAmount)
	}

	trx := &core.SystemTransaction{
		UserId:          id,
		TransactionType: req.GetTransactionType(),
		Amount:          req.GetAmount(),
	}

	err = s.srvc.System.Transaction(ctx, trx)
	if err != nil {
		return nil, status.Error(codes.Internal, errInternalServer)
	}

	return &blnc.SystemTransactionResponse{}, nil
}

func (s *serverAPI) UserTransaction(
	ctx context.Context,
	req *blnc.UserTransactionRequest,
) (*blnc.UserTransactionResponse, error) {
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
