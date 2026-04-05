package service_test

import (
	"context"
	"errors"
	"fmt"
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

func TestSystemService_TransactionTo(t *testing.T) {
	type mockBehavior func(trxTInfo *core.TrxType) []*mock.Call
	const op = "systemsrvc.TransactionTo"

	l := new(logger.LoggerMock)
	system := new(repoMock.System)
	rg := new(trxtyperegistry.RegisterMock)

	sysSrvc := service.NewSystemSrvc(l, system, rg)

	tests := []struct {
		name             string
		input            dto.SystemTrxInput
		rgstyTrxTypeInfo *core.TrxType
		mockBehavior     mockBehavior
		wantErr          error
	}{
		{
			name: "happy path",
			input: dto.SystemTrxInput{
				UserId:  uuid.NewString(),
				Amount:  10000,
				TrxType: 1,
			},
			rgstyTrxTypeInfo: &core.TrxType{Category: "system", Enable: true},
			mockBehavior: func(trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemToType", mock.Anything).Return(trxTInfo, nil)
				c3 := system.On("TransactionTo", mock.Anything, mock.Anything).Return(nil)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: nil,
		},
		{
			name: "uncorrect transaction type",
			input: dto.SystemTrxInput{
				UserId:  uuid.NewString(),
				Amount:  10000,
				TrxType: 0,
			},
			rgstyTrxTypeInfo: &core.TrxType{},
			mockBehavior: func(trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemToType", mock.Anything).Return(trxTInfo, trxtyperegistry.ErrUnknowSysTrxToType)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: core.ErrUnknownTrxType,
		},
		{
			name: "unexpected error from registry",
			input: dto.SystemTrxInput{
				UserId:  uuid.NewString(),
				Amount:  10000,
				TrxType: 1,
			},
			rgstyTrxTypeInfo: &core.TrxType{},
			mockBehavior: func(trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemToType", mock.Anything).Return(trxTInfo, errors.New("unexpected"))
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: errors.New("unexpected"),
		},
		{
			name: "error transaction type is disabled",
			input: dto.SystemTrxInput{
				UserId:  uuid.NewString(),
				Amount:  10000,
				TrxType: 1,
			},
			rgstyTrxTypeInfo: &core.TrxType{Category: "system", Enable: false},
			mockBehavior: func(trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemToType", mock.Anything).Return(trxTInfo, nil)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: core.ErrDisabledType,
		},
		{
			name: "error transaction category not 'system'",
			input: dto.SystemTrxInput{
				UserId:  uuid.NewString(),
				Amount:  10000,
				TrxType: 1,
			},
			rgstyTrxTypeInfo: &core.TrxType{Category: "user", Enable: true},
			mockBehavior: func(trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemToType", mock.Anything).Return(trxTInfo, nil)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: core.ErrInvalidTrxCategory,
		},
		{
			name: "error bad user id",
			input: dto.SystemTrxInput{
				UserId:  "bad uuid",
				Amount:  10000,
				TrxType: 1,
			},
			rgstyTrxTypeInfo: &core.TrxType{Category: "system", Enable: true},
			mockBehavior: func(trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemToType", mock.Anything).Return(trxTInfo, nil)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: core.ErrInvalidUuid,
		},
		{
			name: "error amount equal to 0",
			input: dto.SystemTrxInput{
				UserId:  uuid.NewString(),
				Amount:  0,
				TrxType: 1,
			},
			rgstyTrxTypeInfo: &core.TrxType{Category: "system", Enable: true},
			mockBehavior: func(trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemToType", mock.Anything).Return(trxTInfo, nil)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: core.ErrInvalidAmount,
		},
		{
			name: "error when call db method TransactionTo",
			input: dto.SystemTrxInput{
				UserId:  uuid.NewString(),
				Amount:  10000,
				TrxType: 1,
			},
			rgstyTrxTypeInfo: &core.TrxType{Category: "system", Enable: true},
			mockBehavior: func(trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemToType", mock.Anything).Return(trxTInfo, nil)
				c3 := system.On("TransactionTo", mock.Anything, mock.Anything).Return(fmt.Errorf("err"))
				c4 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3, c4}
			},
			wantErr: fmt.Errorf("err"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := tt.mockBehavior(tt.rgstyTrxTypeInfo)

			err := sysSrvc.TransactionTo(context.Background(), tt.input)

			require.Equal(t, tt.wantErr, err)

			mock.AssertExpectationsForObjects(t, l, system, rg)

			for _, c := range calls {
				c.Unset()
			}
		})
	}
}

func TestSystemService_TransactionFrom(t *testing.T) {
	type mockBehavior func(trxTInfo *core.TrxType) []*mock.Call
	const op = "systemsrvc.TransactionFrom"

	l := new(logger.LoggerMock)
	system := new(repoMock.System)
	rg := new(trxtyperegistry.RegisterMock)

	sysSrvc := service.NewSystemSrvc(l, system, rg)

	tests := []struct {
		name             string
		input            dto.SystemTrxInput
		rgstyTrxTypeInfo *core.TrxType
		mockBehavior     mockBehavior
		wantErr          error
	}{
		{
			name: "happy path",
			input: dto.SystemTrxInput{
				UserId:  uuid.NewString(),
				Amount:  10000,
				TrxType: 1,
			},
			rgstyTrxTypeInfo: &core.TrxType{Category: "system", Enable: true},
			mockBehavior: func(trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemFromType", mock.Anything).Return(trxTInfo, nil)
				c3 := system.On("TransactionFrom", mock.Anything, mock.Anything).Return(nil)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: nil,
		},
		{
			name: "uncorrect transaction type",
			input: dto.SystemTrxInput{
				UserId:  uuid.NewString(),
				Amount:  10000,
				TrxType: 0,
			},
			rgstyTrxTypeInfo: &core.TrxType{},
			mockBehavior: func(trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemFromType", mock.Anything).Return(trxTInfo, trxtyperegistry.ErrUnknowSysTrxFromType)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: core.ErrUnknownTrxType,
		},
		{
			name: "unexpected error from registry",
			input: dto.SystemTrxInput{
				UserId:  uuid.NewString(),
				Amount:  10000,
				TrxType: 1,
			},
			rgstyTrxTypeInfo: &core.TrxType{},
			mockBehavior: func(trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemFromType", mock.Anything).Return(trxTInfo, errors.New("unexpected"))
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: errors.New("unexpected"),
		},
		{
			name: "error transaction type is disabled",
			input: dto.SystemTrxInput{
				UserId:  uuid.NewString(),
				Amount:  10000,
				TrxType: 1,
			},
			rgstyTrxTypeInfo: &core.TrxType{Category: "system", Enable: false},
			mockBehavior: func(trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemFromType", mock.Anything).Return(trxTInfo, nil)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: core.ErrDisabledType,
		},
		{
			name: "error transaction category not 'system'",
			input: dto.SystemTrxInput{
				UserId:  uuid.NewString(),
				Amount:  10000,
				TrxType: 1,
			},
			rgstyTrxTypeInfo: &core.TrxType{Category: "user", Enable: true},
			mockBehavior: func(trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemFromType", mock.Anything).Return(trxTInfo, nil)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: core.ErrInvalidTrxCategory,
		},
		{
			name: "error bad user id",
			input: dto.SystemTrxInput{
				UserId:  "bad uuid",
				Amount:  10000,
				TrxType: 1,
			},
			rgstyTrxTypeInfo: &core.TrxType{Category: "system", Enable: true},
			mockBehavior: func(trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemFromType", mock.Anything).Return(trxTInfo, nil)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: core.ErrInvalidUuid,
		},
		{
			name: "error amount equal to 0",
			input: dto.SystemTrxInput{
				UserId:  uuid.NewString(),
				Amount:  0,
				TrxType: 1,
			},
			rgstyTrxTypeInfo: &core.TrxType{Category: "system", Enable: true},
			mockBehavior: func(trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemFromType", mock.Anything).Return(trxTInfo, nil)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: core.ErrInvalidAmount,
		},
		{
			name: "error when call db method TransactionTo",
			input: dto.SystemTrxInput{
				UserId:  uuid.NewString(),
				Amount:  10000,
				TrxType: 1,
			},
			rgstyTrxTypeInfo: &core.TrxType{Category: "system", Enable: true},
			mockBehavior: func(trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemFromType", mock.Anything).Return(trxTInfo, nil)
				c3 := system.On("TransactionFrom", mock.Anything, mock.Anything).Return(fmt.Errorf("err"))
				c4 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3, c4}
			},
			wantErr: fmt.Errorf("err"),
		},
		{
			name: "error insufficient balance when call db method TransactionFrom",
			input: dto.SystemTrxInput{
				UserId:  uuid.NewString(),
				Amount:  10000,
				TrxType: 1,
			},
			rgstyTrxTypeInfo: &core.TrxType{Category: "system", Enable: true},
			mockBehavior: func(trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemFromType", mock.Anything).Return(trxTInfo, nil)
				c3 := system.On("TransactionFrom", mock.Anything, mock.Anything).Return(postgres.ErrInsuffBalance)
				c4 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3, c4}
			},
			wantErr: core.ErrInsuffBalance,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := tt.mockBehavior(tt.rgstyTrxTypeInfo)

			err := sysSrvc.TransactionFrom(context.Background(), tt.input)

			require.Equal(t, tt.wantErr, err)

			mock.AssertExpectationsForObjects(t, l, system, rg)

			for _, c := range calls {
				c.Unset()
			}
		})
	}
}
