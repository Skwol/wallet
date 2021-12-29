package transaction

import "time"

type TranType string

const (
	TranTypeDeposit  TranType = "deposit"
	TranTypeWithdraw TranType = "withdraw"
	TranTypeTransfer TranType = "transfer"
)

type Transaction struct {
	ID              int64     `json:"id,omitempty"`
	SenderID        int64     `json:"sender_id,omitempty"`
	ReceiverID      int64     `json:"receiver_id,omitempty"`
	Amount          float64   `json:"amount,omitempty"`
	Timestamp       time.Time `json:"timestamp,omitempty"`
	TransactionType TranType  `json:"transaction_type,omitempty"`
}
