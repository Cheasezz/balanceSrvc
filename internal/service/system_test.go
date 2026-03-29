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

func TestSystemService_TransactionTo(t *testing.T) {
	type mockBehavior func(trxT blnc.SystemTrxToType, trxTInfo *core.TrxType) []*mock.Call
	const op = "systemsrvc.TransactionTo"

	l := new(logger.LoggerMock)
	system := new(repo.SystemRepoMock)
	trx := new(repo.TrxRepoMock)
	rp := &repo.Repo{
		System: system,
		Trx:    trx,
	}

	rg := new(trxtyperegistry.RegisterMock)

	sysSrvc := service.NewSystemSrvc(l, rp, rg)

	tests := []struct {
		name             string
		userId           uuid.UUID
		amount           uint64
		trxType          blnc.SystemTrxToType
		rgstyTrxTypeInfo *core.TrxType
		mockBehavior     mockBehavior
		wantErr          error
	}{
		{
			name:             "correct transaction type",
			userId:           uuid.UUID([]byte("1284567891234254")),
			amount:           10000,
			trxType:          blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
			rgstyTrxTypeInfo: &core.TrxType{Category: "system", Enable: true},
			mockBehavior: func(trxT blnc.SystemTrxToType, trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemToType", trxT).Return(trxTInfo, nil)
				c3 := system.On("TransactionTo", mock.Anything, mock.Anything).Return(nil)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: nil,
		},
		{
			name:             "uncorrect transaction type",
			userId:           uuid.UUID([]byte("1234567873234254")),
			amount:           10000,
			trxType:          blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_UNKNOWN,
			rgstyTrxTypeInfo: &core.TrxType{},
			mockBehavior: func(trxT blnc.SystemTrxToType, trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemToType", trxT).Return(trxTInfo, trxtyperegistry.ErrUnknowSysTrxToType)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: fmt.Errorf("%s: %w", op, service.ErrSystemTrxToType),
		},
		{
			name:             "unexpected error from registry",
			userId:           uuid.UUID([]byte("1234567891234254")),
			amount:           10000,
			trxType:          blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_UNKNOWN,
			rgstyTrxTypeInfo: &core.TrxType{},
			mockBehavior: func(trxT blnc.SystemTrxToType, trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemToType", trxT).Return(trxTInfo, errors.New("unexpected"))
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: fmt.Errorf("%s: %w", op, errors.New("unexpected")),
		},
		{
			name:             "error transaction type is disabled",
			userId:           uuid.UUID([]byte("1234567891234985")),
			amount:           10000,
			trxType:          blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
			rgstyTrxTypeInfo: &core.TrxType{Category: "system", Enable: false},
			mockBehavior: func(trxT blnc.SystemTrxToType, trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemToType", trxT).Return(trxTInfo, nil)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: fmt.Errorf("%s: %w", op, core.ErrDisabledType),
		},
		{
			name:             "error transaction category not 'system'",
			userId:           uuid.UUID([]byte("1234567891234321")),
			amount:           10000,
			trxType:          blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
			rgstyTrxTypeInfo: &core.TrxType{Category: "user", Enable: true},
			mockBehavior: func(trxT blnc.SystemTrxToType, trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemToType", trxT).Return(trxTInfo, nil)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: fmt.Errorf("%s: %w", op, core.ErrInvalidTrxCategory),
		},
		{
			name:             "error bad user id (uuid.Nil)",
			userId:           uuid.Nil,
			amount:           10000,
			trxType:          blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
			rgstyTrxTypeInfo: &core.TrxType{Category: "system", Enable: true},
			mockBehavior: func(trxT blnc.SystemTrxToType, trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemToType", trxT).Return(trxTInfo, nil)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: fmt.Errorf("%s: %w", op, core.ErrInvalidUserId),
		},
		{
			name:             "error amount equal to 0",
			userId:           uuid.UUID([]byte("1234567891234123")),
			amount:           0,
			trxType:          blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
			rgstyTrxTypeInfo: &core.TrxType{Category: "system", Enable: true},
			mockBehavior: func(trxT blnc.SystemTrxToType, trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemToType", trxT).Return(trxTInfo, nil)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: fmt.Errorf("%s: %w", op, core.ErrInvalidAmount),
		},
		{
			name:             "error when call db method TransactionTo",
			userId:           uuid.UUID([]byte("12345678912345375")),
			amount:           10000,
			trxType:          blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
			rgstyTrxTypeInfo: &core.TrxType{Category: "system", Enable: true},
			mockBehavior: func(trxT blnc.SystemTrxToType, trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemToType", trxT).Return(trxTInfo, nil)
				c3 := system.On("TransactionTo", mock.Anything, mock.Anything).Return(fmt.Errorf("err"))
				c4 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3, c4}
			},
			wantErr: fmt.Errorf("%s: %w", op, fmt.Errorf("err")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := tt.mockBehavior(tt.trxType, tt.rgstyTrxTypeInfo)

			err := sysSrvc.TransactionTo(context.Background(), tt.userId, tt.amount, tt.trxType)

			require.Equal(t, tt.wantErr, err)

			mock.AssertExpectationsForObjects(t, l, system, trx, rg)

			for _, c := range calls {
				c.Unset()
			}
		})
	}
}

func TestSystemService_TransactionFrom(t *testing.T) {
	type mockBehavior func(trxT blnc.SystemTrxFromType, trxTInfo *core.TrxType) []*mock.Call
	const op = "systemsrvc.TransactionFrom"

	l := new(logger.LoggerMock)
	system := new(repo.SystemRepoMock)
	trx := new(repo.TrxRepoMock)
	rp := &repo.Repo{
		System: system,
		Trx:    trx,
	}

	rg := new(trxtyperegistry.RegisterMock)

	sysSrvc := service.NewSystemSrvc(l, rp, rg)

	tests := []struct {
		name             string
		userId           uuid.UUID
		amount           uint64
		trxType          blnc.SystemTrxFromType
		rgstyTrxTypeInfo *core.TrxType
		mockBehavior     mockBehavior
		wantErr          error
	}{
		{
			name:             "happy path",
			userId:           uuid.UUID([]byte("1284567891234254")),
			amount:           10000,
			trxType:          blnc.SystemTrxFromType_SYSTEM_TRX_FROM_TYPE_WITHDRAWAL,
			rgstyTrxTypeInfo: &core.TrxType{Category: "system", Enable: true},
			mockBehavior: func(trxT blnc.SystemTrxFromType, trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemFromType", trxT).Return(trxTInfo, nil)
				c3 := system.On("TransactionFrom", mock.Anything, mock.Anything).Return(nil)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: nil,
		},
		{
			name:             "uncorrect transaction type",
			userId:           uuid.UUID([]byte("1234567873234254")),
			amount:           10000,
			trxType:          blnc.SystemTrxFromType_SYSTEM_TRX_FROM_TYPE_WITHDRAWAL,
			rgstyTrxTypeInfo: &core.TrxType{},
			mockBehavior: func(trxT blnc.SystemTrxFromType, trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemFromType", trxT).Return(trxTInfo, trxtyperegistry.ErrUnknowSysTrxFromType)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: fmt.Errorf("%s: %w", op, service.ErrSystemTrxFromType),
		},
		{
			name:             "unexpected error from registry",
			userId:           uuid.UUID([]byte("1234567891234254")),
			amount:           10000,
			trxType:          blnc.SystemTrxFromType_SYSTEM_TRX_FROM_TYPE_WITHDRAWAL,
			rgstyTrxTypeInfo: &core.TrxType{},
			mockBehavior: func(trxT blnc.SystemTrxFromType, trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemFromType", trxT).Return(trxTInfo, errors.New("unexpected"))
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: fmt.Errorf("%s: %w", op, errors.New("unexpected")),
		},
		{
			name:             "error transaction type is disabled",
			userId:           uuid.UUID([]byte("1234567891234985")),
			amount:           10000,
			trxType:          blnc.SystemTrxFromType_SYSTEM_TRX_FROM_TYPE_WITHDRAWAL,
			rgstyTrxTypeInfo: &core.TrxType{Category: "system", Enable: false},
			mockBehavior: func(trxT blnc.SystemTrxFromType, trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemFromType", trxT).Return(trxTInfo, nil)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: fmt.Errorf("%s: %w", op, core.ErrDisabledType),
		},
		{
			name:             "error transaction category not 'system'",
			userId:           uuid.UUID([]byte("1234567891234321")),
			amount:           10000,
			trxType:          blnc.SystemTrxFromType_SYSTEM_TRX_FROM_TYPE_WITHDRAWAL,
			rgstyTrxTypeInfo: &core.TrxType{Category: "user", Enable: true},
			mockBehavior: func(trxT blnc.SystemTrxFromType, trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemFromType", trxT).Return(trxTInfo, nil)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: fmt.Errorf("%s: %w", op, core.ErrInvalidTrxCategory),
		},
		{
			name:             "error bad user id (uuid.Nil)",
			userId:           uuid.Nil,
			amount:           10000,
			trxType:          blnc.SystemTrxFromType_SYSTEM_TRX_FROM_TYPE_WITHDRAWAL,
			rgstyTrxTypeInfo: &core.TrxType{Category: "system", Enable: true},
			mockBehavior: func(trxT blnc.SystemTrxFromType, trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemFromType", trxT).Return(trxTInfo, nil)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: fmt.Errorf("%s: %w", op, core.ErrInvalidUserId),
		},
		{
			name:             "error amount equal to 0",
			userId:           uuid.UUID([]byte("1234567891234123")),
			amount:           0,
			trxType:          blnc.SystemTrxFromType_SYSTEM_TRX_FROM_TYPE_WITHDRAWAL,
			rgstyTrxTypeInfo: &core.TrxType{Category: "system", Enable: true},
			mockBehavior: func(trxT blnc.SystemTrxFromType, trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemFromType", trxT).Return(trxTInfo, nil)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: fmt.Errorf("%s: %w", op, core.ErrInvalidAmount),
		},
		{
			name:             "error when call db method TransactionTo",
			userId:           uuid.UUID([]byte("12345678912345375")),
			amount:           10000,
			trxType:          blnc.SystemTrxFromType_SYSTEM_TRX_FROM_TYPE_WITHDRAWAL,
			rgstyTrxTypeInfo: &core.TrxType{Category: "system", Enable: true},
			mockBehavior: func(trxT blnc.SystemTrxFromType, trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemFromType", trxT).Return(trxTInfo, nil)
				c3 := system.On("TransactionFrom", mock.Anything, mock.Anything).Return(fmt.Errorf("err"))
				c4 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3, c4}
			},
			wantErr: fmt.Errorf("%s: %w", op, fmt.Errorf("err")),
		},
		{
			name:             "error insufficient balance when call db method TransactionFrom",
			userId:           uuid.UUID([]byte("12345678912345375")),
			amount:           10000,
			trxType:          blnc.SystemTrxFromType_SYSTEM_TRX_FROM_TYPE_WITHDRAWAL,
			rgstyTrxTypeInfo: &core.TrxType{Category: "system", Enable: true},
			mockBehavior: func(trxT blnc.SystemTrxFromType, trxTInfo *core.TrxType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemFromType", trxT).Return(trxTInfo, nil)
				c3 := system.On("TransactionFrom", mock.Anything, mock.Anything).Return(repo.ErrInsuffBalance)
				c4 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3, c4}
			},
			wantErr: fmt.Errorf("%s: %w", op, service.ErrInsuffBalance),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := tt.mockBehavior(tt.trxType, tt.rgstyTrxTypeInfo)

			err := sysSrvc.TransactionFrom(context.Background(), tt.userId, tt.amount, tt.trxType)

			require.Equal(t, tt.wantErr, err)

			mock.AssertExpectationsForObjects(t, l, system, trx, rg)

			for _, c := range calls {
				c.Unset()
			}
		})
	}
}
