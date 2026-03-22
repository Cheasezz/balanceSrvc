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
	type mockBehavior func(trxT blnc.SystemTrxToType) []*mock.Call
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
		name         string
		userId       uuid.UUID
		amount       int64
		trxType      blnc.SystemTrxToType
		mockBehavior mockBehavior
		wantErr      error
	}{
		{
			name:    "correct transaction type",
			userId:  uuid.Max,
			amount:  10000,
			trxType: blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
			mockBehavior: func(trxT blnc.SystemTrxToType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemToType", trxT).Return(&core.TrxType{}, nil)
				c3 := system.On("TransactionTo", mock.Anything, mock.Anything).Return(nil)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: nil,
		},
		{
			name:    "uncorrect transaction type",
			userId:  uuid.Max,
			amount:  10000,
			trxType: blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_UNKNOWN,
			mockBehavior: func(trxT blnc.SystemTrxToType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemToType", trxT).Return(&core.TrxType{}, trxtyperegistry.ErrUnknowSysTrxToType)
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: fmt.Errorf("%s: %w", op, service.ErrSystemTrxToType),
		},
		{
			name:    "unexpected error from registry",
			userId:  uuid.Max,
			amount:  10000,
			trxType: blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_UNKNOWN,
			mockBehavior: func(trxT blnc.SystemTrxToType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemToType", trxT).Return(&core.TrxType{}, errors.New("unexpected"))
				c3 := l.On("Error", mock.Anything, mock.Anything, mock.Anything)
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: fmt.Errorf("%s: %w", op, errors.New("unexpected")),
		},
		{
			name:    "error when call db method TransactionTo",
			userId:  uuid.Max,
			amount:  10000,
			trxType: blnc.SystemTrxToType_SYSTEM_TRX_TO_TYPE_DEPOSIT,
			mockBehavior: func(trxT blnc.SystemTrxToType) []*mock.Call {
				c1 := l.On("With", "op", op).Return(l)
				c2 := rg.On("SystemToType", trxT).Return(&core.TrxType{}, nil)
				c3 := system.On("TransactionTo", mock.Anything, mock.Anything).Return(fmt.Errorf("err"))
				return []*mock.Call{c1, c2, c3}
			},
			wantErr: fmt.Errorf("%s: %w", op, fmt.Errorf("err")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := tt.mockBehavior(tt.trxType)

			err := sysSrvc.TransactionTo(context.Background(), tt.userId, tt.amount, tt.trxType)

			require.Equal(t, tt.wantErr, err)

			mock.AssertExpectationsForObjects(t, l, system, trx, rg)

			for _, c := range calls {
				c.Unset()
			}
		})
	}
}
