package transfer

import "time"

type TransferDTO struct {
	CreateTransferDTO
}

type CreateTransferDTO struct {
	ID        int64     `json:"id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
	Sender    WalletDTO `json:"sender"`
	Receiver  WalletDTO `json:"receiver"`
}

func (d CreateTransferDTO) toModel() *Transfer {
	return &Transfer{
		Amount:    d.Amount,
		Timestamp: d.Timestamp,
		Sender:    d.Sender.toModel(),
		Receiver:  d.Receiver.toModel(),
	}
}

type WalletDTO struct {
	ID      int64   `json:"id"`
	Balance float64 `json:"balance"`
}

func (d WalletDTO) toModel() *Wallet {
	return &Wallet{
		ID:      d.ID,
		Balance: d.Balance,
	}
}
