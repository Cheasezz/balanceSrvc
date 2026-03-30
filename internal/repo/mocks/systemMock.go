package repoMock

import (
	"context"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/stretchr/testify/mock"
)

type System struct {
	mock.Mock
}

func (m *System) TransactionTo(ctx context.Context, trx *core.Transaction) error {
	args := m.Called(ctx, trx)
	return args.Error(0)
}

func (m *System) TransactionFrom(ctx context.Context, trx *core.Transaction) error {
	args := m.Called(ctx, trx)
	return args.Error(0)
}
