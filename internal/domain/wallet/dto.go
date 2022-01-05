package wallet

import (
	"time"
)

func walletsToDTO(wallets []*Wallet) []WalletDTO {
	result := make([]WalletDTO, len(wallets))
	for i, wallet := range wallets {
		result[i] = wallet.toDTO()
	}
	return result
}

type WalletDTO struct {
	ID                  int64            `json:"id"`
	Name                string           `json:"name"`
	Balance             float64          `json:"balance"`
	TransactionsToApply []TransactionDTO `json:"-"`
	Transactions        []TransactionDTO `json:"transactions,omitempty"`
}

func (d WalletDTO) toModel() *Wallet {
	return &Wallet{
		ID:      d.ID,
		Name:    d.Name,
		Balance: d.Balance,
	}
}

type CreateWalletDTO struct {
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

type UpdateWalletDTO struct {
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

type TransactionDTO struct {
	ID         int64     `json:"id,omitempty"`
	SenderID   int64     `json:"sender_id,omitempty"`
	ReceiverID int64     `json:"receiver_id,omitempty"`
	Amount     float64   `json:"amount,omitempty"`
	Timestamp  time.Time `json:"timestamp,omitempty"`
	Type       TranType  `json:"type,omitempty"`
}
