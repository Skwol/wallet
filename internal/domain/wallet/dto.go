package wallet

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
