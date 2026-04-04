package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Cheasezz/balanceSrvc/internal/adapter/postgres"
	repoMock "github.com/Cheasezz/balanceSrvc/internal/adapter/postgres/mocks"
	trxtyperegistry "github.com/Cheasezz/balanceSrvc/internal/adapter/trxTypeRegistry"
	"github.com/Cheasezz/balanceSrvc/internal/core"
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
	system := new(repoMock.System)
	user := new(repoMock.User)
	trx := new(repoMock.Trx)
	rp := &postgres.Postgres{
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
			wantErr: service.ErrUsrTrxType,
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
			wantErr: errors.New("unexpected"),
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
			wantErr: service.ErrUserTrxTypeDisabled,
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
			wantErr: service.ErrSameIds,
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
			wantErr: core.ErrInvalidTrxCategory,
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
				c3 := user.On("TransactionToUser", mock.Anything, mock.Anything).Return(postgres.ErrInsuffBalance)
				c4 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3, c4}
			},
			wantErr: service.ErrInsuffBalance,
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
			wantErr: errors.New("err"),
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

func TestUserService_Balance(t *testing.T) {
	type mockBehavior func() []*mock.Call
	const op = "usersrvc.Balance"

	l := new(logger.LoggerMock)
	system := new(repoMock.System)
	user := new(repoMock.User)
	trx := new(repoMock.Trx)
	rp := &postgres.Postgres{
		System: system,
		User:   user,
		Trx:    trx,
	}

	rg := new(trxtyperegistry.RegisterMock)

	usrSrvc := service.NewUserSrvc(l, rp, rg)
	tests := []struct {
		name         string
		userId       uuid.UUID
		mockBehavior mockBehavior
		wantErr      error
	}{
		{
			name:   "happy path",
			userId: uuid.New(),
			mockBehavior: func() []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := user.On("Balance", mock.Anything, mock.Anything).Return(10, nil)
				return []*mock.Call{c1, c2}
			},
			wantErr: nil,
		},
		{
			name:   "db error id not found",
			userId: uuid.New(),
			mockBehavior: func() []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := user.On("Balance", mock.Anything, mock.Anything).Return(0, postgres.ErrIdNotfound)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: service.ErrIdNotfound,
		},
		{
			name:   "unexpected error from postgres layer",
			userId: uuid.New(),
			mockBehavior: func() []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := user.On("Balance", mock.Anything, mock.Anything).Return(0, errors.New("err"))
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: errors.New("err"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := tt.mockBehavior()

			_, err := usrSrvc.Balance(context.Background(), tt.userId)

			require.Equal(t, tt.wantErr, err)

			mock.AssertExpectationsForObjects(t, l, system, trx, rg)

			for _, c := range calls {
				c.Unset()
			}
		})
	}
}
