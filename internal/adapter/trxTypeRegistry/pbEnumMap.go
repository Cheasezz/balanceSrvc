package trxtyperegistry

import (
	"fmt"
	"strings"

	"github.com/Cheasezz/balanceSrvc/internal/core"
)

const (
	errTrxNotFound = "transaction type from enum protobuff not found in DB"
)

func buildEnumMap(
	enumNames map[int32]string,
	dbTrxTypes map[string]*core.TrxType,
	prefix string,
) (map[int32]*core.TrxType, error) {
	const op = "trxtyperegistry.buildEnumMap"

	res := make(map[int32]*core.TrxType)

	for val, name := range enumNames {
		if val == 0 {
			continue
		}

		code := strings.ToLower(strings.TrimPrefix(name, prefix))

		info, ok := dbTrxTypes[code]
		if !ok {
			return nil, fmt.Errorf("op=%s, err=%s, trxType=%s", op, errTrxNotFound, code)
		}
		res[val] = info
	}
	return res, nil
}
