package trxtyperegistry

import (
	"github.com/Cheasezz/balanceSrvc/internal/core"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
	"github.com/stretchr/testify/mock"
)

type RegisterMock struct {
	mock.Mock
}

func (m *RegisterMock) SystemToType(t blnc.SystemTrxToType) (*core.TrxType, error) {
	agrgs := m.Called(t)
	return agrgs.Get(0).(*core.TrxType), agrgs.Error(1)
}

func (m *RegisterMock) SystemFromType(t blnc.SystemTrxFromType) (*core.TrxType, error) {
	agrgs := m.Called(t)
	return agrgs.Get(0).(*core.TrxType), agrgs.Error(1)
}

func (m *RegisterMock) UserType(t blnc.UserTrxType) (*core.TrxType, error) {
	agrgs := m.Called(t)
	return agrgs.Get(0).(*core.TrxType), agrgs.Error(1)
}
