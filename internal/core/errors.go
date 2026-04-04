package core

import "errors"

var (
	ErrInvalidUuid        = errors.New("field user_id must be valid uuid")
	ErrIdNotfound         = errors.New("id not found")
	ErrUnknownTrxType     = errors.New("unknown transaction type")
	ErrDisabledType       = errors.New("this type is disabled")
	ErrInvalidAmount      = errors.New("invalid amount value, must be uint and not equal to 0")
	ErrInsuffBalance      = errors.New("insufficient balance")
	ErrInvalidTrxCategory = errors.New("this type of transaction not in current catagory")
	ErrSameIds            = errors.New("Ids must be not equal")
	ErrInternalServer     = errors.New("something went wrong on server")
)
