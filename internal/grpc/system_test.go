package grpcSrv_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	grpcSrv "github.com/Cheasezz/balanceSrvc/internal/grpc"
	"github.com/Cheasezz/balanceSrvc/internal/service"
	srvcMock "github.com/Cheasezz/balanceSrvc/internal/service/mocks"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestSystemHandler_TransactionTo(t *testing.T) {
	type mockBehavior func() []*mock.Call

	sysSrvc := new(srvcMock.System)
	s := &service.Service{
		System: sysSrvc,
	}

	handlers := grpcSrv.ServerAPI{
		blnc.UnimplementedBalanceServer{},
		s,
	}

	tests := []struct {
		name         string
		req          *blnc.SystemTrxToRequest
		mockBehavior mockBehavior
		wantResp     *blnc.SystemTrxResponse
		wantErr      error
	}{
		{
			name: "happy path",
			req: &blnc.SystemTrxToRequest{
				UserId:        "37166f7a-f430-49e9-8306-8fba9fbf4311",
				SystemTrxType: blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
				Amount:        10000,
			},
			mockBehavior: func() []*mock.Call {
				c1 := sysSrvc.On("TransactionTo", mock.Anything, mock.Anything).Return(nil)
				return []*mock.Call{c1}
			},
			wantResp: &blnc.SystemTrxResponse{},
			wantErr:  nil,
		},
		{
			name: "error bad uuid",
			req: &blnc.SystemTrxToRequest{
				UserId:        "baaaad uuid",
				SystemTrxType: blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
				Amount:        10000,
			},
			mockBehavior: func() []*mock.Call {
				c1 := sysSrvc.On("TransactionTo", mock.Anything, mock.Anything).
					Return(core.ErrInvalidUuid)
				return []*mock.Call{c1}
			},
			wantResp: nil,
			wantErr:  status.Error(codes.InvalidArgument, core.ErrInvalidUuid.Error()),
		},
		{
			name: "error zero amount",
			req: &blnc.SystemTrxToRequest{
				UserId:        uuid.NewString(),
				SystemTrxType: blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
				Amount:        0,
			},
			mockBehavior: func() []*mock.Call {
				c1 := sysSrvc.On("TransactionTo", mock.Anything, mock.Anything).Return(core.ErrInvalidAmount)
				return []*mock.Call{c1}
			},
			wantResp: nil,
			wantErr:  status.Error(codes.InvalidArgument, core.ErrInvalidAmount.Error()),
		},
		{
			name: "error service check uncorrect transaction type",
			req: &blnc.SystemTrxToRequest{
				UserId:        uuid.NewString(),
				SystemTrxType: blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_UNKNOWN,
				Amount:        10000,
			},
			mockBehavior: func() []*mock.Call {
				c1 := sysSrvc.On(
					"TransactionTo", mock.Anything, mock.Anything,
				).Return(core.ErrUnknownTrxType)
				return []*mock.Call{c1}
			},
			wantResp: nil,
			wantErr:  status.Error(codes.InvalidArgument, core.ErrUnknownTrxType.Error()),
		},
		{
			name: "error service check disabled transaction type",
			req: &blnc.SystemTrxToRequest{
				UserId:        uuid.NewString(),
				SystemTrxType: blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
				Amount:        10000,
			},
			mockBehavior: func() []*mock.Call {
				c1 := sysSrvc.On(
					"TransactionTo", mock.Anything, mock.Anything,
				).Return(core.ErrDisabledType)
				return []*mock.Call{c1}
			},
			wantResp: nil,
			wantErr:  status.Error(codes.InvalidArgument, core.ErrDisabledType.Error()),
		},
		{
			name: "unexpected error when check transaction type in service",
			req: &blnc.SystemTrxToRequest{
				UserId:        "37166f7a-f430-49e9-8306-8fba9fbf4311",
				SystemTrxType: blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
				Amount:        10000,
			},
			mockBehavior: func() []*mock.Call {
				c1 := sysSrvc.On(
					"TransactionTo", mock.Anything, mock.Anything,
				).Return(errors.New("unexpected error"))
				return []*mock.Call{c1}
			},
			wantResp: nil,
			wantErr:  status.Error(codes.Internal, core.ErrInternalServer.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := tt.mockBehavior()

			resp, err := handlers.SystemTransactionTo(context.Background(), tt.req)
			require.Equal(t, tt.wantResp, resp)
			require.Equal(t, tt.wantErr, err)

			mock.AssertExpectationsForObjects(t, sysSrvc)

			for _, c := range calls {
				c.Unset()
			}
		})
	}
}

func TestSystemHandler_TransactionFrom(t *testing.T) {
	type mockBehavior func() []*mock.Call

	sysSrvc := new(srvcMock.System)
	s := &service.Service{
		System: sysSrvc,
	}

	handlers := grpcSrv.ServerAPI{
		blnc.UnimplementedBalanceServer{},
		s,
	}

	tests := []struct {
		name         string
		req          *blnc.SystemTrxFromRequest
		mockBehavior mockBehavior
		wantResp     *blnc.SystemTrxResponse
		wantErr      error
	}{
		{
			name: "happy path",
			req: &blnc.SystemTrxFromRequest{
				UserId:        uuid.NewString(),
				SystemTrxType: blnc.SystemTrxFromType_SYSTEM_TRX_FROM_TYPE_WITHDRAWAL,
				Amount:        10000,
			},
			mockBehavior: func() []*mock.Call {
				c1 := sysSrvc.On("TransactionFrom", mock.Anything, mock.Anything).Return(nil)
				return []*mock.Call{c1}
			},
			wantResp: &blnc.SystemTrxResponse{},
			wantErr:  nil,
		},
		{
			name: "error bad uuid",
			req: &blnc.SystemTrxFromRequest{
				UserId:        "baaaad uuid",
				SystemTrxType: blnc.SystemTrxFromType_SYSTEM_TRX_FROM_TYPE_WITHDRAWAL,
				Amount:        10000,
			},
			mockBehavior: func() []*mock.Call {
				c1 := sysSrvc.On("TransactionFrom", mock.Anything, mock.Anything).Return(core.ErrInvalidUuid)
				return []*mock.Call{c1}
			},
			wantResp: nil,
			wantErr:  status.Error(codes.InvalidArgument, core.ErrInvalidUuid.Error()),
		},
		{
			name: "error zero amount",
			req: &blnc.SystemTrxFromRequest{
				UserId:        uuid.NewString(),
				SystemTrxType: blnc.SystemTrxFromType_SYSTEM_TRX_FROM_TYPE_WITHDRAWAL,
				Amount:        0,
			},
			mockBehavior: func() []*mock.Call {
				c1 := sysSrvc.On("TransactionFrom", mock.Anything, mock.Anything).Return(core.ErrInvalidAmount)
				return []*mock.Call{c1}
			},
			wantResp: nil,
			wantErr:  status.Error(codes.InvalidArgument, core.ErrInvalidAmount.Error()),
		},
		{
			name: "error service check uncorrect transaction type",
			req: &blnc.SystemTrxFromRequest{
				UserId:        uuid.NewString(),
				SystemTrxType: blnc.SystemTrxFromType_SYSTEM_TRX_FROM_TYPE_WITHDRAWAL,
				Amount:        10000,
			},
			mockBehavior: func() []*mock.Call {
				c1 := sysSrvc.On(
					"TransactionFrom", mock.Anything, mock.Anything,
				).Return(core.ErrUnknownTrxType)
				return []*mock.Call{c1}
			},
			wantResp: nil,
			wantErr:  status.Error(codes.InvalidArgument, core.ErrUnknownTrxType.Error()),
		},
		{
			name: "error service check disabled transaction type",
			req: &blnc.SystemTrxFromRequest{
				UserId:        uuid.NewString(),
				SystemTrxType: blnc.SystemTrxFromType_SYSTEM_TRX_FROM_TYPE_WITHDRAWAL,
				Amount:        10000,
			},
			mockBehavior: func() []*mock.Call {
				c1 := sysSrvc.On(
					"TransactionFrom", mock.Anything, mock.Anything,
				).Return(core.ErrDisabledType)
				return []*mock.Call{c1}
			},
			wantResp: nil,
			wantErr:  status.Error(codes.InvalidArgument, core.ErrDisabledType.Error()),
		},
		{
			name: "error service insufficient balance",
			req: &blnc.SystemTrxFromRequest{
				UserId:        uuid.NewString(),
				SystemTrxType: blnc.SystemTrxFromType_SYSTEM_TRX_FROM_TYPE_WITHDRAWAL,
				Amount:        10000,
			},
			mockBehavior: func() []*mock.Call {
				c1 := sysSrvc.On(
					"TransactionFrom", mock.Anything, mock.Anything, mock.Anything,
				).Return(core.ErrInsuffBalance)
				return []*mock.Call{c1}
			},
			wantResp: nil,
			wantErr:  status.Error(codes.InvalidArgument, core.ErrInsuffBalance.Error()),
		},
		{
			name: "unexpected error when check transaction type in service",
			req: &blnc.SystemTrxFromRequest{
				UserId:        uuid.NewString(),
				SystemTrxType: blnc.SystemTrxFromType_SYSTEM_TRX_FROM_TYPE_WITHDRAWAL,
				Amount:        10000,
			},
			mockBehavior: func() []*mock.Call {
				c1 := sysSrvc.On(
					"TransactionFrom", mock.Anything, mock.Anything, mock.Anything,
				).Return(errors.New("unexpected error"))
				return []*mock.Call{c1}
			},
			wantResp: nil,
			wantErr:  status.Error(codes.Internal, core.ErrInternalServer.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := tt.mockBehavior()

			resp, err := handlers.SystemTransactionFrom(context.Background(), tt.req)
			assert.Equal(t, tt.wantResp, resp)
			assert.Equal(t, tt.wantErr, err)

			mock.AssertExpectationsForObjects(t, sysSrvc)

			for _, c := range calls {
				c.Unset()
			}
		})
	}
}
