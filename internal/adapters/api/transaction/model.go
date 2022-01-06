package transaction

import (
	"fmt"
	"time"

	"github.com/skwol/wallet/internal/domain/transaction"
)

var csvHeaders = []string{"Transaction ID", "Sender ID", "Receiver ID", "Amount", "Timestamp", "Type"}

func newTransaction(dto transaction.TransactionDTO) Transaction {
	return Transaction{
		ID:         dto.ID,
		SenderID:   dto.SenderID,
		ReceiverID: dto.ReceiverID,
		Amount:     dto.Amount,
		Timestamp:  dto.Timestamp,
		Type:       string(dto.Type),
	}
}

type Transaction struct {
	ID         int64     `json:"id"`
	SenderID   int64     `json:"sender_id"`
	ReceiverID int64     `json:"receiver_id"`
	Amount     float64   `json:"amount"`
	Timestamp  time.Time `json:"timestamp"`
	Type       string    `json:"type"`
}

func (t Transaction) toCsv() []string {
	return []string{fmt.Sprintf("%d", t.ID), fmt.Sprintf("%d", t.SenderID), fmt.Sprintf("%d", t.ReceiverID), fmt.Sprintf("%f", t.Amount), t.Timestamp.Format("Mon, 02 Jan 2006 15:04:05 -0700"), t.Type}
}

type Filter struct {
	SenderIDs   []int64          `json:"sender_ids"`
	ReceiverIDs []int64          `json:"receiver_ids"`
	Amount      FloatRangeFilter `json:"amount"`
	Timestamp   DateRangeFilter  `json:"timestamp"`
	Types       []string         `json:"types"`
}

func (f Filter) toFilterRequest() transaction.FilterTransactionsDTO {
	return transaction.FilterTransactionsDTO{
		SenderIDs:   f.SenderIDs,
		ReceiverIDs: f.ReceiverIDs,
		Amount:      f.Amount.toRequest(),
		Timestamp:   f.Timestamp.toRequest(),
		Types:       f.Types,
	}
}

type FloatRangeFilter struct {
	From float64 `json:"from"`
	To   float64 `json:"to"`
}

func (f FloatRangeFilter) toRequest() transaction.FloatRangeFilter {
	return transaction.FloatRangeFilter{
		From: f.From,
		To:   f.To,
	}
}

type DateRangeFilter struct {
	From time.Time `json:"from"`
	To   time.Time `json:"to"`
}

func (f DateRangeFilter) toRequest() transaction.DateRangeFilter {
	return transaction.DateRangeFilter{
		From: f.From,
		To:   f.To,
	}
}
