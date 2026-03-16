package trxtyperegistry

import (
	"fmt"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
)

const (
	errUnknowSysTrxType = "unknow system transaction type"
	errUnknowUsrTrxType = "unknow user transaction type"
)

const (
	systemTypePrefix = "SYSTEM_TRX_TYPE_"
	userTypePrefix   = "USER_TRX_TYPE_"
)

type Registry struct {
	system map[blnc.SystemTrxType]*core.TrxType
	user   map[blnc.UserTrxType]*core.TrxType
}

func New(trxTypes map[string]*core.TrxType) (*Registry, error) {
	system, err := buildEnumMap[blnc.SystemTrxType](
		blnc.SystemTrxType_name,
		trxTypes,
		systemTypePrefix,
	)
	if err != nil {
		return nil, err
	}

	user, err := buildEnumMap[blnc.UserTrxType](
		blnc.UserTrxType_name,
		trxTypes,
		userTypePrefix,
	)
	if err != nil {
		return nil, err
	}

	return &Registry{system: system, user: user}, nil
}

func (r *Registry) SystemType(t blnc.SystemTrxType) (*core.TrxType, error) {
	const op = "trxtyperegistry.SystemType"

	info, ok := r.system[t]
	if !ok {
		return nil, fmt.Errorf("op=%s, err=%s", op, errUnknowSysTrxType)
	}

	return info, nil
}

func (r *Registry) UserType(t blnc.UserTrxType) (*core.TrxType, error) {
	const op = "trxtyperegistry.UserType"

	info, ok := r.user[t]
	if !ok {
		return nil, fmt.Errorf("op=%s, err=%s", op, errUnknowUsrTrxType)
	}

	return info, nil
}
