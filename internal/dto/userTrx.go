package dto

type UserTrxInput struct {
	Sender    string
	Resipient string
	Amount    uint64
	TrxType   int32
}
