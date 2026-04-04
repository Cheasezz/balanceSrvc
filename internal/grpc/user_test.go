package grpcSrv_test

import (
	"context"
	"errors"
	"testing"

	grpcSrv "github.com/Cheasezz/balanceSrvc/internal/grpc"
	"github.com/Cheasezz/balanceSrvc/internal/service"
	srvcMock "github.com/Cheasezz/balanceSrvc/internal/service/mocks"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUserHandler_TransactionToUser(t *testing.T) {
	type mockBehavior func() []*mock.Call

	usrSrvc := new(srvcMock.User)
	s := &service.Service{
		User: usrSrvc,
	}

	handlers := grpcSrv.ServerAPI{
		blnc.UnimplementedBalanceServer{},
		s,
	}
	u := uuid.NewString()

	tests := []struct {
		name         string
		req          *blnc.UserTrxRequest
		mockBehavior mockBehavior
		wantResp     *blnc.UserTrxResponse
		wantErr      error
	}{
		{
			name: "happy path",
			req: &blnc.UserTrxRequest{
				SenderId:    uuid.NewString(),
				ResipientId: uuid.NewString(),
				UserTrxType: blnc.UserTrxType_USER_TRX_TYPE_TRANSFER,
				Amount:      10000,
			},
			mockBehavior: func() []*mock.Call {
				c1 := usrSrvc.On(
					"TransactionToUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
				).Return(nil)
				return []*mock.Call{c1}
			},

			wantResp: &blnc.UserTrxResponse{},
			wantErr:  nil,
		},
		{
			name: "error sender bad uuid",
			req: &blnc.UserTrxRequest{
				SenderId:    "baaaad uuid",
				ResipientId: uuid.NewString(),
				UserTrxType: blnc.UserTrxType_USER_TRX_TYPE_TRANSFER,
				Amount:      10000,
			},
			mockBehavior: func() []*mock.Call { return []*mock.Call{} },
			wantResp:     nil,
			wantErr:      status.Error(codes.InvalidArgument, grpcSrv.ErrInvalidUuid.Error()),
		},
		{
			name: "error resipient bad uuid",
			req: &blnc.UserTrxRequest{
				SenderId:    uuid.NewString(),
				ResipientId: "baaaad uuid",
				UserTrxType: blnc.UserTrxType_USER_TRX_TYPE_TRANSFER,
				Amount:      10000,
			},
			mockBehavior: func() []*mock.Call { return []*mock.Call{} },
			wantResp:     nil,
			wantErr:      status.Error(codes.InvalidArgument, grpcSrv.ErrInvalidUuid.Error()),
		},
		{
			name: "error zero amount",
			req: &blnc.UserTrxRequest{
				SenderId:    uuid.NewString(),
				ResipientId: uuid.NewString(),
				UserTrxType: blnc.UserTrxType_USER_TRX_TYPE_TRANSFER,
				Amount:      0,
			},
			mockBehavior: func() []*mock.Call { return []*mock.Call{} },
			wantResp:     nil,
			wantErr:      status.Error(codes.InvalidArgument, grpcSrv.ErrInvalidAmount.Error()),
		},
		{
			name: "error service check uncorrect transaction type",
			req: &blnc.UserTrxRequest{
				SenderId:    uuid.NewString(),
				ResipientId: uuid.NewString(),
				UserTrxType: blnc.UserTrxType_USER_TRX_TYPE_UNKNOWN,
				Amount:      10000,
			},
			mockBehavior: func() []*mock.Call {
				c1 := usrSrvc.On(
					"TransactionToUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
				).Return(service.ErrUsrTrxType)
				return []*mock.Call{c1}
			},
			wantResp: nil,
			wantErr:  status.Error(codes.InvalidArgument, service.ErrUsrTrxType.Error()),
		},
		{
			name: "error service check disabled transaction type",
			req: &blnc.UserTrxRequest{
				SenderId:    uuid.NewString(),
				ResipientId: uuid.NewString(),
				UserTrxType: blnc.UserTrxType_USER_TRX_TYPE_TRANSFER,
				Amount:      10000,
			},
			mockBehavior: func() []*mock.Call {
				c1 := usrSrvc.On(
					"TransactionToUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
				).Return(service.ErrUserTrxTypeDisabled)
				return []*mock.Call{c1}
			},
			wantResp: nil,
			wantErr:  status.Error(codes.InvalidArgument, service.ErrUserTrxTypeDisabled.Error()),
		},
		{
			name: "error insufficient balance",
			req: &blnc.UserTrxRequest{
				SenderId:    uuid.NewString(),
				ResipientId: uuid.NewString(),
				UserTrxType: blnc.UserTrxType_USER_TRX_TYPE_TRANSFER,
				Amount:      10000,
			},
			mockBehavior: func() []*mock.Call {
				c1 := usrSrvc.On(
					"TransactionToUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
				).Return(service.ErrInsuffBalance)
				return []*mock.Call{c1}
			},
			wantResp: nil,
			wantErr:  status.Error(codes.InvalidArgument, service.ErrInsuffBalance.Error()),
		},
		{
			name: "error same ids",
			req: &blnc.UserTrxRequest{
				SenderId:    u,
				ResipientId: u,
				UserTrxType: blnc.UserTrxType_USER_TRX_TYPE_TRANSFER,
				Amount:      10000,
			},
			mockBehavior: func() []*mock.Call {
				c1 := usrSrvc.On(
					"TransactionToUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
				).Return(service.ErrSameIds)
				return []*mock.Call{c1}
			},
			wantResp: nil,
			wantErr:  status.Error(codes.InvalidArgument, service.ErrSameIds.Error()),
		},
		{
			name: "unexpected error service",
			req: &blnc.UserTrxRequest{
				SenderId:    u,
				ResipientId: u,
				UserTrxType: blnc.UserTrxType_USER_TRX_TYPE_TRANSFER,
				Amount:      10000,
			},
			mockBehavior: func() []*mock.Call {
				c1 := usrSrvc.On(
					"TransactionToUser", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
				).Return(errors.New("err"))
				return []*mock.Call{c1}
			},
			wantResp: nil,
			wantErr:  status.Error(codes.Internal, grpcSrv.ErrInternalServer.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := tt.mockBehavior()

			resp, err := handlers.UserTransaction(context.Background(), tt.req)
			require.Equal(t, tt.wantResp, resp)
			require.Equal(t, tt.wantErr, err)

			mock.AssertExpectationsForObjects(t, usrSrvc)

			for _, c := range calls {
				c.Unset()
			}
		})
	}
}

func TestUserHandler_UserBalance(t *testing.T) {
	type mockBehavior func() []*mock.Call

	usrSrvc := new(srvcMock.User)
	s := &service.Service{
		User: usrSrvc,
	}

	handlers := grpcSrv.ServerAPI{
		blnc.UnimplementedBalanceServer{},
		s,
	}

	tests := []struct {
		name         string
		req          *blnc.BalanceRequest
		mockBehavior mockBehavior
		wantResp     *blnc.BalanceResponse
		wantErr      error
	}{
		{
			name: "happy path",
			req:  &blnc.BalanceRequest{UserId: uuid.NewString()},
			mockBehavior: func() []*mock.Call {
				c1 := usrSrvc.On("Balance", mock.Anything, mock.Anything).Return(100000, nil)
				return []*mock.Call{c1}
			},
			wantResp: &blnc.BalanceResponse{Balance: 100000},
			wantErr:  nil,
		},
		{
			name: "error bad uuid",
			req:  &blnc.BalanceRequest{UserId: "bad uuid"},
			mockBehavior: func() []*mock.Call {
				return []*mock.Call{}
			},
			wantResp: nil,
			wantErr:  status.Error(codes.InvalidArgument, grpcSrv.ErrInvalidUuid.Error()),
		},
		{
			name: "error uuid not found",
			req:  &blnc.BalanceRequest{UserId: uuid.NewString()},
			mockBehavior: func() []*mock.Call {
				c1 := usrSrvc.On("Balance", mock.Anything, mock.Anything).Return(0, service.ErrIdNotfound)
				return []*mock.Call{c1}
			},
			wantResp: nil,
			wantErr:  status.Error(codes.NotFound, service.ErrIdNotfound.Error()),
		},
		{
			name: "unexpected error from service layer",
			req:  &blnc.BalanceRequest{UserId: uuid.NewString()},
			mockBehavior: func() []*mock.Call {
				c1 := usrSrvc.On("Balance", mock.Anything, mock.Anything).Return(0, errors.New("err"))
				return []*mock.Call{c1}
			},
			wantResp: nil,
			wantErr:  status.Error(codes.Internal, grpcSrv.ErrInternalServer.Error()),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := tt.mockBehavior()

			resp, err := handlers.UserBalance(context.Background(), tt.req)
			require.Equal(t, tt.wantResp, resp)
			require.Equal(t, tt.wantErr, err)

			mock.AssertExpectationsForObjects(t, usrSrvc)

			for _, c := range calls {
				c.Unset()
			}
		})
	}
}
