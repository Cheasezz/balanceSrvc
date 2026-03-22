package repo

import (
	"context"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/stretchr/testify/mock"
)

type TrxRepoMock struct {
	mock.Mock
}

func (m *TrxRepoMock) GetAllTypesInfo(ctx context.Context) (map[string]*core.TrxType, error) {
	args := m.Called(ctx)
	return args.Get(0).(map[string]*core.TrxType), args.Error(1)
}
