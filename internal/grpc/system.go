package grpcSrv

import (
	"context"

	"github.com/Cheasezz/balanceSrvc/internal/dto"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
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
		return nil, toStatus(err)
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
		return nil, toStatus(err)
	}

	return &blnc.SystemTrxResponse{}, nil
}
