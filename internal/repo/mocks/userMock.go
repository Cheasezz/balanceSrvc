package repoMock

import (
	"context"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/stretchr/testify/mock"
)

type User struct {
	mock.Mock
}

func (m *User) TransactionToUser(ctx context.Context, trx *core.Transaction) error {
	args := m.Called(ctx, trx)
	return args.Error(0)
}
