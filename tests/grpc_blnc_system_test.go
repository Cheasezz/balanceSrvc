package tests

import (
	"context"
	"testing"

	grpcHndlrs "github.com/Cheasezz/balanceSrvc/internal/grpc"
	"github.com/Cheasezz/balanceSrvc/internal/service"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
	testsuite "github.com/Cheasezz/balanceSrvc/tests/suite"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGrpcBalance_SystemTransactionTo(t *testing.T) {
	t.Parallel()

	suit := testsuite.New(t)

	tests := []struct {
		name    string
		req     *blnc.SystemTrxToRequest
		wantErr error
	}{
		{
			name: "happy path",
			req: &blnc.SystemTrxToRequest{
				UserId:        uuid.NewString(),
				SystemTrxType: blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
				Amount:        10000,
			},
			wantErr: nil,
		},
		{
			name: "error bad userId",
			req: &blnc.SystemTrxToRequest{
				UserId:        "bad uuid",
				SystemTrxType: blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
				Amount:        10000,
			},
			wantErr: status.Error(codes.InvalidArgument, grpcHndlrs.ErrInvalidUuid.Error()),
		},
		{
			name: "error zero amount",
			req: &blnc.SystemTrxToRequest{
				UserId:        uuid.NewString(),
				SystemTrxType: blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
				Amount:        0,
			},
			wantErr: status.Error(codes.InvalidArgument, grpcHndlrs.ErrInvalidAmount.Error()),
		},
		{
			name: "error invalid transaction type",
			req: &blnc.SystemTrxToRequest{
				UserId:        uuid.NewString(),
				SystemTrxType: blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_UNKNOWN,
				Amount:        10000,
			},
			wantErr: status.Error(codes.InvalidArgument, service.ErrSystemTrxToType.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel()

			ctx, cancelCtx := context.WithTimeout(context.Background(), suit.CtxTimeout)
			defer cancelCtx()

			_, err := suit.BalanceClient.SystemTransactionTo(ctx, tt.req)

			require.ErrorIs(t, err, tt.wantErr)
		})
	}
}
