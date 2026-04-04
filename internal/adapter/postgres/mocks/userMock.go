package repoMock

import (
	"context"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type User struct {
	mock.Mock
}

func (m *User) TransactionToUser(ctx context.Context, trx *core.Transaction) error {
	args := m.Called(ctx, trx)
	return args.Error(0)
}

func (m *User) Balance(ctx context.Context, userId uuid.UUID) (int64, error) {
	args := m.Called(ctx, userId)
	return int64(args.Int(0)), args.Error(1)
}
