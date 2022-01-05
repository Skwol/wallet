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
	ID                  int64
	Name                string
	Balance             float64
	TransactionsToApply []TransactionDTO
	Transactions        []TransactionDTO
}

func (d WalletDTO) toModel() Wallet {
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
