package wallet

import "github.com/skwol/wallet/internal/domain/transaction"

type WalletDTO struct {
	ID                  int64                         `json:"id"`
	Name                string                        `json:"name"`
	Balance             float64                       `json:"balance"`
	TransactionsToApply []*transaction.TransactionDTO `json:"-"`
}

func (d WalletDTO) toModel() *Wallet {
	return &Wallet{
		ID:      d.ID,
		Name:    d.Name,
		Balance: d.Balance,
	}
}

func walletsToDTO(wallets []*Wallet) []*WalletDTO {
	result := make([]*WalletDTO, len(wallets))
	for i, wallet := range wallets {
		result[i] = wallet.toDTO()
	}
	return result
}

type CreateWalletDTO struct {
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}

type UpdateWalletDTO struct {
	Name    string  `json:"name"`
	Balance float64 `json:"balance"`
}
