package srvcMock

import (
	"context"

	"github.com/Cheasezz/balanceSrvc/internal/dto"
	"github.com/stretchr/testify/mock"
)

type User struct {
	mock.Mock
}

func (m *User) TransactionToUser(ctx context.Context, input dto.UserTrxInput) error {
	args := m.Called(ctx, input)
	return args.Error(0)
}

func (m *User) Balance(ctx context.Context, userId string) (uint64, error) {
	args := m.Called(ctx, userId)
	return uint64(args.Int(0)), args.Error(1)
}
