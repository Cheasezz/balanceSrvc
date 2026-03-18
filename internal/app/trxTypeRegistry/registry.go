package trxtyperegistry

import (
	"errors"
	"fmt"

	"github.com/Cheasezz/balanceSrvc/internal/core"
	blnc "github.com/Cheasezz/balanceSrvc/protos/gen"
)

var (
	ErrUnknowSysTrxToType   = errors.New("unknow system transaction(to) type")
	ErrUnknowSysTrxFromType = errors.New("unknow system transaction(from) type")
	ErrUnknowUsrTrxType     = errors.New("unknow user transaction type")
)

const (
	systemToTypePrefix   = "SYSTEM_TRX_TO_TYPE_"
	systemFromTypePrefix = "SYSTEM_TRX_FROM_TYPE_"
	userTypePrefix       = "USER_TRX_TYPE_"
)

type Registry struct {
	systemTo   map[blnc.SystemTrxToType]*core.TrxType
	systemFrom map[blnc.SystemTrxFromType]*core.TrxType
	user       map[blnc.UserTrxType]*core.TrxType
}

func New(trxTypes map[string]*core.TrxType) (*Registry, error) {
	systemTo, err := buildEnumMap[blnc.SystemTrxToType](
		blnc.SystemTrxToType_name,
		trxTypes,
		systemToTypePrefix,
	)
	if err != nil {
		return nil, err
	}

	systemFrom, err := buildEnumMap[blnc.SystemTrxFromType](
		blnc.SystemTrxFromType_name,
		trxTypes,
		systemFromTypePrefix,
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

	return &Registry{systemTo: systemTo, systemFrom: systemFrom, user: user}, nil
}

func (r *Registry) SystemToType(t blnc.SystemTrxToType) (*core.TrxType, error) {
	const op = "trxtyperegistry.SystemToType"

	info, ok := r.systemTo[t]
	if !ok {
		return nil, fmt.Errorf("%s: %w", op, ErrUnknowSysTrxToType)
	}

	return info, nil
}

func (r *Registry) SystemFromType(t blnc.SystemTrxFromType) (*core.TrxType, error) {
	const op = "trxtyperegistry.SystemFromType"

	info, ok := r.systemFrom[t]
	if !ok {
		return nil, fmt.Errorf("op=%s, err=%s", op, ErrUnknowSysTrxFromType)
	}

	return info, nil
}

func (r *Registry) UserType(t blnc.UserTrxType) (*core.TrxType, error) {
	const op = "trxtyperegistry.UserType"

	info, ok := r.user[t]
	if !ok {
		return nil, fmt.Errorf("op=%s, err=%s", op, ErrUnknowUsrTrxType)
	}

	return info, nil
}
