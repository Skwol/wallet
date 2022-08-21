package wallet

import (
	"time"

	"github.com/skwol/wallet/internal/domain/wallet"
)

func newWallet(dto wallet.DTO) Wallet {
	transactions := make([]Transaction, 0, len(dto.Transactions))
	for _, t := range dto.Transactions {
		transactions = append(transactions, newTransaction(t))
	}
	return Wallet{
		ID:           dto.ID,
		Name:         dto.Name,
		Balance:      dto.Balance,
		Transactions: transactions,
	}
}

type Wallet struct {
	ID           int64         `json:"id"`
	Name         string        `json:"name"`
	Balance      float64       `json:"balance"`
	Transactions []Transaction `json:"transactions,omitempty"`
}

func (w Wallet) toCreateRequest() wallet.CreateWalletDTO {
	return wallet.CreateWalletDTO{
		Name:    w.Name,
		Balance: w.Balance,
	}
}

func (w Wallet) toUpdateRequest() wallet.UpdateWalletDTO {
	return wallet.UpdateWalletDTO{
		CreateWalletDTO: w.toCreateRequest(),
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

func newTransaction(dto wallet.TransactionDTO) Transaction {
	return Transaction{
		ID:         dto.ID,
		SenderID:   dto.SenderID,
		ReceiverID: dto.ReceiverID,
		Amount:     dto.Amount,
		Timestamp:  dto.Timestamp,
		Type:       string(dto.Type),
	}
}
