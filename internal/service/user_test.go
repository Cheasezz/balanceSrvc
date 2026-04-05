package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Cheasezz/balanceSrvc/internal/adapter/postgres"
	repoMock "github.com/Cheasezz/balanceSrvc/internal/adapter/postgres/mocks"
	trxtyperegistry "github.com/Cheasezz/balanceSrvc/internal/adapter/trxTypeRegistry"
	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/Cheasezz/balanceSrvc/internal/dto"
	"github.com/Cheasezz/balanceSrvc/internal/service"
	"github.com/Cheasezz/balanceSrvc/pkg/logger"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUserService_TransactionToUser(t *testing.T) {
	type mockBehavior func(trxTInfo *core.TrxType) []*mock.Call
	const op = "usersrvc.TransactionToUser"

	l := new(logger.LoggerMock)
	user := new(repoMock.User)
	rg := new(trxtyperegistry.RegisterMock)
	usrSrvc := service.NewUserSrvc(l, user, rg)

	u := uuid.NewString()
	tests := []struct {
		name             string
		input            dto.UserTrxInput
		rgstyTrxTypeInfo *core.TrxType
		mockBehavior     mockBehavior
		wantErr          error
	}{
		{
			name: "happy path",
			input: dto.UserTrxInput{
				Sender:    uuid.NewString(),
				Resipient: uuid.NewString(),
				Amount:    10000,
				TrxType:   1,
			},
			rgstyTrxTypeInfo: &core.TrxType{Category: "user", Enable: true},
			mockBehavior: func(trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("UserType", mock.Anything).Return(trxTInfo, nil)
				c3 := user.On("TransactionToUser", mock.Anything, mock.Anything).Return(nil)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: nil,
		},
		{
			name: "uncorrect transaction type",
			input: dto.UserTrxInput{
				Sender:    uuid.NewString(),
				Resipient: uuid.NewString(),
				Amount:    10000,
				TrxType:   0,
			},
			rgstyTrxTypeInfo: &core.TrxType{},
			mockBehavior: func(trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("UserType", mock.Anything).Return(trxTInfo, trxtyperegistry.ErrUnknowUsrTrxType)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: core.ErrUnknownTrxType,
		},
		{
			name: "unexpected error from registry",
			input: dto.UserTrxInput{
				Sender:    uuid.NewString(),
				Resipient: uuid.NewString(),
				Amount:    10000,
				TrxType:   1,
			},
			rgstyTrxTypeInfo: &core.TrxType{},
			mockBehavior: func(trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("UserType", mock.Anything).Return(trxTInfo, errors.New("unexpected"))
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: errors.New("unexpected"),
		},
		{
			name: "error transaction type is disabled",
			input: dto.UserTrxInput{
				Sender:    uuid.NewString(),
				Resipient: uuid.NewString(),
				Amount:    10000,
				TrxType:   1,
			},
			rgstyTrxTypeInfo: &core.TrxType{Category: "user", Enable: false},
			mockBehavior: func(trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("UserType", mock.Anything).Return(trxTInfo, nil)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: core.ErrDisabledType,
		},
		{
			name: "error same ids",
			input: dto.UserTrxInput{
				Sender:    u,
				Resipient: u,
				Amount:    10000,
				TrxType:   1,
			},
			rgstyTrxTypeInfo: &core.TrxType{Category: "user", Enable: true},
			mockBehavior: func(trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("UserType", mock.Anything).Return(trxTInfo, nil)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: core.ErrSameIds,
		},
		{
			name: "error transaction category not 'user'",
			input: dto.UserTrxInput{
				Sender:    uuid.NewString(),
				Resipient: uuid.NewString(),
				Amount:    10000,
				TrxType:   1,
			},
			rgstyTrxTypeInfo: &core.TrxType{Category: "system", Enable: true},
			mockBehavior: func(trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("UserType", mock.Anything).Return(trxTInfo, nil)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: core.ErrInvalidTrxCategory,
		},
		{
			name: "error insufficient balance db method TransactionToUser",
			input: dto.UserTrxInput{
				Sender:    uuid.NewString(),
				Resipient: uuid.NewString(),
				Amount:    10000,
				TrxType:   1,
			},
			rgstyTrxTypeInfo: &core.TrxType{Category: "user", Enable: true},
			mockBehavior: func(trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("UserType", mock.Anything).Return(trxTInfo, nil)
				c3 := user.On("TransactionToUser", mock.Anything, mock.Anything).Return(postgres.ErrInsuffBalance)
				c4 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3, c4}
			},
			wantErr: core.ErrInsuffBalance,
		},
		{
			name: "error when call db method TransactionToUser",
			input: dto.UserTrxInput{
				Sender:    uuid.NewString(),
				Resipient: uuid.NewString(),
				Amount:    10000,
				TrxType:   1,
			},
			rgstyTrxTypeInfo: &core.TrxType{Category: "user", Enable: true},
			mockBehavior: func(trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("UserType", mock.Anything).Return(trxTInfo, nil)
				c3 := user.On("TransactionToUser", mock.Anything, mock.Anything).Return(errors.New("err"))
				c4 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3, c4}
			},
			wantErr: errors.New("err"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := tt.mockBehavior(tt.rgstyTrxTypeInfo)

			err := usrSrvc.TransactionToUser(context.Background(), tt.input)

			require.Equal(t, tt.wantErr, err)

			mock.AssertExpectationsForObjects(t, l, rg)

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
	user := new(repoMock.User)
	rg := new(trxtyperegistry.RegisterMock)
	usrSrvc := service.NewUserSrvc(l, user, rg)

	tests := []struct {
		name         string
		userId       string
		mockBehavior mockBehavior
		wantErr      error
	}{
		{
			name:   "happy path",
			userId: uuid.NewString(),
			mockBehavior: func() []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := user.On("Balance", mock.Anything, mock.Anything).Return(10, nil)
				return []*mock.Call{c1, c2}
			},
			wantErr: nil,
		},
		{
			name:   "bad uuid",
			userId: "bad uuid",
			mockBehavior: func() []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := user.On("Balance", mock.Anything, mock.Anything).Return(0, core.ErrInvalidUuid)
				return []*mock.Call{c1, c2}
			},
			wantErr: core.ErrInvalidUuid,
		},
		{
			name:   "db error id not found",
			userId: uuid.NewString(),
			mockBehavior: func() []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := user.On("Balance", mock.Anything, mock.Anything).Return(0, postgres.ErrIdNotfound)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: core.ErrIdNotfound,
		},
		{
			name:   "unexpected error from postgres layer",
			userId: uuid.NewString(),
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

			mock.AssertExpectationsForObjects(t, l, rg)

			for _, c := range calls {
				c.Unset()
			}
		})
	}
}
