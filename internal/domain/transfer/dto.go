package transfer

import "time"

type DTO struct {
	ID int64
	CreateTransferDTO
}

type CreateTransferDTO struct {
	Amount    float64
	Timestamp time.Time
	Sender    WalletDTO
	Receiver  WalletDTO
}

func (d CreateTransferDTO) validate() error {
	if d.Sender.ID == 0 {
		return ErrMissingSender
	}
	if d.Receiver.ID == 0 {
		return ErrMissingReceiver
	}
	if d.Receiver.ID == d.Sender.ID {
		return ErrSameSenderAndReceiver
	}
	if d.Amount <= 0 {
		return ErrNonPositiveAmount
	}
	if d.Sender.Balance-d.Amount < 0 {
		return ErrNotEnoughMoney
	}
	return nil
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
	return Wallet(d)
}
