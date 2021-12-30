package transaction

import "time"

type TransactionDTO struct {
	ID         int64     `json:"id,omitempty"`
	SenderID   int64     `json:"sender_id,omitempty"`
	ReceiverID int64     `json:"receiver_id,omitempty"`
	Amount     float64   `json:"amount,omitempty"`
	Timestamp  time.Time `json:"timestamp,omitempty"`
	Type       TranType  `json:"type,omitempty"`
}

type CreateTransactionDTO struct {
	SenderID   int64     `json:"sender_id"`
	ReceiverID int64     `json:"receiver_id"`
	Amount     float64   `json:"amount"`
	Timestamp  time.Time `json:"timestamp"`
	Type       string    `json:"type"`
}
