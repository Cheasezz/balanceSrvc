package grpcHndlrs

import (
	"context"

	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
	"google.golang.org/grpc"
)

type serverAPI struct {
	blnc.UnimplementedBalanceServer
}

func Register(gRPC *grpc.Server) {
	blnc.RegisterBalanceServer(gRPC, &serverAPI{})
}

func (s *serverAPI) SystemTransaction(
	ctx context.Context,
	req *blnc.SystemTransactionRequest,
) (*blnc.SystemTransactionResponse, error) {
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
