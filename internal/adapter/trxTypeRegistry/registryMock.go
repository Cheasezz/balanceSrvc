package trxtyperegistry

import (
	"github.com/Cheasezz/balanceSrvc/internal/core"
	"github.com/stretchr/testify/mock"
)

type RegisterMock struct {
	mock.Mock
}

func (m *RegisterMock) SystemToType(t int32) (*core.TrxType, error) {
	agrgs := m.Called(t)
	return agrgs.Get(0).(*core.TrxType), agrgs.Error(1)
}

func (m *RegisterMock) SystemFromType(t int32) (*core.TrxType, error) {
	agrgs := m.Called(t)
	return agrgs.Get(0).(*core.TrxType), agrgs.Error(1)
}

func (m *RegisterMock) UserType(t int32) (*core.TrxType, error) {
	agrgs := m.Called(t)
	return agrgs.Get(0).(*core.TrxType), agrgs.Error(1)
}
