package transaction

import "time"

type CreateTransactionDTO struct {
	SenderID        int64     `json:"sender_id"`
	ReceiverID      int64     `json:"receiver_id"`
	Amount          float64   `json:"amount"`
	Timestamp       time.Time `json:"timestamp"`
	TransactionType string    `json:"transaction_type"`
}
