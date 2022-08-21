package wallet

import (
	"time"
)

type DTO struct {
	ID                  int64
	Name                string
	Balance             float64
	TransactionsToApply []TransactionDTO
	Transactions        []TransactionDTO
}

func (d DTO) toModel() Wallet {
	return Wallet{
		ID:      d.ID,
		Name:    d.Name,
		Balance: d.Balance,
	}
}

type CreateWalletDTO struct {
	Name    string
	Balance float64
}

type UpdateWalletDTO struct {
	CreateWalletDTO
}

type TransactionDTO struct {
	ID         int64
	SenderID   int64
	ReceiverID int64
	Amount     float64
	Timestamp  time.Time
	Type       TranType
}
