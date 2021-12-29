package wallet

type WalletDTO struct {
	ID        int64   `json:"id"`
	Name      string  `json:"name"`
	AccountID int64   `json:"account_id"`
	Balance   float64 `json:"balance"`
}

func walletToDTO(wallet *Wallet) *WalletDTO {
	return &WalletDTO{
		ID:        wallet.ID,
		Name:      wallet.Name,
		AccountID: wallet.AccountID,
		Balance:   wallet.Balance,
	}
}

func walletsToDTO(wallets []*Wallet) []*WalletDTO {
	result := make([]*WalletDTO, len(wallets))
	for i, wallet := range wallets {
		result[i] = walletToDTO(wallet)
	}
	return result
}

type CreateWalletDTO struct {
	Name      string  `json:"name"`
	AccountID int64   `json:"account_id"`
	Balance   float64 `json:"balance"`
}

type UpdateWalletDTO struct {
	Name      string  `json:"name"`
	AccountID int64   `json:"account_id"`
	Balance   float64 `json:"balance"`
}
