package srvcMock

import (
	"context"

	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type User struct {
	mock.Mock
}

func (m *User) TransactionToUser(
	ctx context.Context,
	sender,
	resipient uuid.UUID,
	amount uint64,
	trxType blnc.UserTrxType,
) error {
	args := m.Called(ctx, sender, resipient, amount, trxType)
	return args.Error(0)
}
