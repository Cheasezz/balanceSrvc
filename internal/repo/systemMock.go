package repo

import (
	"context"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/stretchr/testify/mock"
)

type SystemRepoMock struct {
	mock.Mock
}

func (m *SystemRepoMock) TransactionTo(ctx context.Context, trx *core.Transaction) error {
	args := m.Called(ctx, trx)
	return args.Error(0)
}
