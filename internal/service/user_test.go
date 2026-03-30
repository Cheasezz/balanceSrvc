package service_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	trxtyperegistry "github.com/Cheasezz/balanceSrvc/internal/app/trxTypeRegistry"
	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/Cheasezz/balanceSrvc/internal/repo"
	"github.com/Cheasezz/balanceSrvc/internal/service"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserService_TransactionToUser(t *testing.T) {
	type mockBehavior func(trxT blnc.UserTrxType, trxTInfo *core.TrxType) []*mock.Call
	const op = "usersrvc.TransactionToUser"

	l := new(logger.LoggerMock)
	system := new(repo.SystemRepoMock)
	user := new(repo.UserRepoMock)
	trx := new(repo.TrxRepoMock)
	rp := &repo.Repo{
		System: system,
		User:   user,
		Trx:    trx,
	}

	rg := new(trxtyperegistry.RegisterMock)

	usrSrvc := service.NewUserSrvc(l, rp, rg)
	u := uuid.New()
	tests := []struct {
		name             string
		sender           uuid.UUID
		resipient        uuid.UUID
		amount           uint64
		trxType          blnc.UserTrxType
		rgstyTrxTypeInfo *core.TrxType
		mockBehavior     mockBehavior
		wantErr          error
	}{
		{
			name:             "happy path",
			sender:           uuid.New(),
			resipient:        uuid.New(),
			amount:           10000,
			trxType:          blnc.UserTrxType_USER_TRX_TYPE_TRANSFER,
			rgstyTrxTypeInfo: &core.TrxType{Category: "user", Enable: true},
			mockBehavior: func(trxT blnc.UserTrxType, trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("UserType", trxT).Return(trxTInfo, nil)
				c3 := user.On("TransactionToUser", mock.Anything, mock.Anything).Return(nil)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: nil,
		},
		{
			name:             "uncorrect transaction type",
			sender:           uuid.New(),
			resipient:        uuid.New(),
			amount:           10000,
			trxType:          blnc.UserTrxType_USER_TRX_TYPE_UNKNOWN,
			rgstyTrxTypeInfo: &core.TrxType{},
			mockBehavior: func(trxT blnc.UserTrxType, trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("UserType", trxT).Return(trxTInfo, trxtyperegistry.ErrUnknowUsrTrxType)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: fmt.Errorf("%s: %w", op, service.ErrUsrTrxType),
		},
		{
			name:             "unexpected error from registry",
			sender:           uuid.New(),
			resipient:        uuid.New(),
			amount:           10000,
			trxType:          blnc.UserTrxType_USER_TRX_TYPE_TRANSFER,
			rgstyTrxTypeInfo: &core.TrxType{},
			mockBehavior: func(trxT blnc.UserTrxType, trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("UserType", trxT).Return(trxTInfo, errors.New("unexpected"))
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: fmt.Errorf("%s: %w", op, errors.New("unexpected")),
		},
		{
			name:             "error transaction type is disabled",
			sender:           uuid.New(),
			resipient:        uuid.New(),
			amount:           10000,
			trxType:          blnc.UserTrxType_USER_TRX_TYPE_TRANSFER,
			rgstyTrxTypeInfo: &core.TrxType{Category: "user", Enable: false},
			mockBehavior: func(trxT blnc.UserTrxType, trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("UserType", trxT).Return(trxTInfo, nil)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: fmt.Errorf("%s: %w", op, service.ErrUserTrxTypeDisabled),
		},
		{
			name:             "error same ids",
			sender:           u,
			resipient:        u,
			amount:           10000,
			trxType:          blnc.UserTrxType_USER_TRX_TYPE_TRANSFER,
			rgstyTrxTypeInfo: &core.TrxType{Category: "user", Enable: true},
			mockBehavior: func(trxT blnc.UserTrxType, trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("UserType", trxT).Return(trxTInfo, nil)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: fmt.Errorf("%s: %w", op, service.ErrSameIds),
		},
		{
			name:             "error transaction category not 'user'",
			sender:           uuid.New(),
			resipient:        uuid.New(),
			amount:           10000,
			trxType:          blnc.UserTrxType_USER_TRX_TYPE_TRANSFER,
			rgstyTrxTypeInfo: &core.TrxType{Category: "system", Enable: true},
			mockBehavior: func(trxT blnc.UserTrxType, trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("UserType", trxT).Return(trxTInfo, nil)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: fmt.Errorf("%s: %w", op, core.ErrInvalidTrxCategory),
		},
		{
			name:             "error insufficient balance db method TransactionToUser",
			sender:           uuid.New(),
			resipient:        uuid.New(),
			amount:           10000,
			trxType:          blnc.UserTrxType_USER_TRX_TYPE_TRANSFER,
			rgstyTrxTypeInfo: &core.TrxType{Category: "user", Enable: true},
			mockBehavior: func(trxT blnc.UserTrxType, trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("UserType", trxT).Return(trxTInfo, nil)
				c3 := user.On("TransactionToUser", mock.Anything, mock.Anything).Return(repo.ErrInsuffBalance)
				c4 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3, c4}
			},
			wantErr: fmt.Errorf("%s: %w", op, service.ErrInsuffBalance),
		},
		{
			name:             "error when call db method TransactionToUser",
			sender:           uuid.New(),
			resipient:        uuid.New(),
			amount:           10000,
			trxType:          blnc.UserTrxType_USER_TRX_TYPE_TRANSFER,
			rgstyTrxTypeInfo: &core.TrxType{Category: "user", Enable: true},
			mockBehavior: func(trxT blnc.UserTrxType, trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("UserType", trxT).Return(trxTInfo, nil)
				c3 := user.On("TransactionToUser", mock.Anything, mock.Anything).Return(errors.New("err"))
				c4 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3, c4}
			},
			wantErr: fmt.Errorf("%s: %w", op, errors.New("err")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := tt.mockBehavior(tt.trxType, tt.rgstyTrxTypeInfo)

			err := usrSrvc.TransactionToUser(context.Background(), tt.sender, tt.resipient, tt.amount, tt.trxType)

			require.Equal(t, tt.wantErr, err)

			mock.AssertExpectationsForObjects(t, l, system, trx, rg)

			for _, c := range calls {
				c.Unset()
			}
		})
	}
}
