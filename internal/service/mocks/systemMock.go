package srvcMock

import (
	"context"

	"github.com/Cheasezz/balanceSrvc/internal/dto"
	"github.com/stretchr/testify/mock"
)

type System struct {
	mock.Mock
}

func (m *System) TransactionTo(ctx context.Context, input dto.SystemTrxInput) error {

	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *System) TransactionFrom(ctx context.Context, input dto.SystemTrxInput) error {

	args := m.Called(ctx, input)
	return args.Error(0)
}
