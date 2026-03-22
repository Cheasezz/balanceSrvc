package service

import (
	"context"

	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type SysSrvcMock struct {
	mock.Mock
}

func (m *SysSrvcMock) TransactionTo(
	ctx context.Context,
	userId uuid.UUID,
	amount uint64,
	trxType blnc.SystemTrxToType,
) error {

	args := m.Called(ctx, userId, amount, trxType)
	return args.Error(0)
}
