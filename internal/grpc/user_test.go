package grpcHndlrs_test

import (
	"context"
	"errors"
	"testing"

	grpcHndlrs "github.com/Cheasezz/balanceSrvc/internal/grpc"
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

	handlers := grpcHndlrs.ServerAPI{
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
			wantErr:      status.Error(codes.InvalidArgument, grpcHndlrs.ErrInvalidUuid.Error()),
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
			wantErr:      status.Error(codes.InvalidArgument, grpcHndlrs.ErrInvalidUuid.Error()),
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
			wantErr:      status.Error(codes.InvalidArgument, grpcHndlrs.ErrInvalidAmount.Error()),
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
			wantErr:  status.Error(codes.Internal, grpcHndlrs.ErrInternalServer.Error()),
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
