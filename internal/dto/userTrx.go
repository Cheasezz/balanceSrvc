package dto

type UserTrxInput struct {
	IdempotencyKey string
	Sender         string
	Resipient      string
	Amount         uint64
	TrxType        int32
}
