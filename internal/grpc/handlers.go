package grpcHndlrs

import (
	"context"

	"github.com/Cheasezz/balanceSrvc/internal/service"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
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
	l    logger.Logger
}

func Register(gRPC *grpc.Server, l logger.Logger, s *service.Service) {
	blnc.RegisterBalanceServer(gRPC, &serverAPI{srvc: s, l: l})
}

func (s *serverAPI) SystemTransaction(
	ctx context.Context,
	req *blnc.SystemTrxRequest,
) (*blnc.SystemTrxResponse, error) {

	const op = "grpcHndlrs.SystemTransaction"
	log := s.l.With("op", op)

	id, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, errInvalidUuid)
	}

	if req.GetAmount() == 0 {
		return nil, status.Error(codes.InvalidArgument, errInvalidAmount)
	}

	err = s.srvc.System.Transaction(ctx, id, req.GetAmount(), req.SystemTrxType)
	if err != nil {
		log.Error(err.Error())
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
