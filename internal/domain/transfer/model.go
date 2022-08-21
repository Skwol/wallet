package transfer

import (
	"fmt"
	"time"
)

type Transfer struct {
	Amount    float64
	Timestamp time.Time
	Sender    Wallet
	Receiver  Wallet
}

func (t *Transfer) toDTO() *DTO {
	return &DTO{
		CreateTransferDTO: CreateTransferDTO{
			Amount:    t.Amount,
			Timestamp: t.Timestamp,
			Receiver:  t.Receiver.toDTO(),
			Sender:    t.Sender.toDTO(),
		},
	}
}

type Wallet struct {
	ID      int64
	Balance float64
}

func (w *Wallet) toDTO() WalletDTO {
	return WalletDTO{
		ID:      w.ID,
		Balance: w.Balance,
	}
}

func createTransfer(dto *CreateTransferDTO, timestamp time.Time) (*Transfer, error) {
	if dto.Receiver.ID == 0 || dto.Sender.ID == 0 {
		return nil, fmt.Errorf("missing sender or receiver")
	}
	if dto.Receiver.ID == dto.Sender.ID {
		return nil, fmt.Errorf("transfer can not be performed when sender and receiver is the same wallet")
	}
	if dto.Amount <= 0 {
		return nil, fmt.Errorf("amount should be greater then 0")
	}
	if dto.Sender.Balance-dto.Amount < 0 {
		return nil, fmt.Errorf("sender does not have enough 'money' for transfer")
	}

	dto.Sender.Balance -= dto.Amount
	dto.Receiver.Balance += dto.Amount
	dto.Timestamp = timestamp
	return dto.toModel(), nil
}
