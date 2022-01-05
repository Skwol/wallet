package transfer

import "time"

type TransferDTO struct {
	ID int64
	CreateTransferDTO
}

type CreateTransferDTO struct {
	Amount    float64
	Timestamp time.Time
	Sender    WalletDTO
	Receiver  WalletDTO
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
	ID      int64
	Balance float64
}

func (d WalletDTO) toModel() Wallet {
	return Wallet{
		ID:      d.ID,
		Balance: d.Balance,
	}
}
