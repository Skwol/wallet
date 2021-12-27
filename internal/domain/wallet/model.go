package wallet

type Wallet struct {
	ID        int64   `json:"id,omitempty"`
	Name      string  `json:"name,omitempty"`
	AccountID int64   `json:"account_id,omitempty"`
	Balance   float64 `json:"balance,omitempty"`
}
