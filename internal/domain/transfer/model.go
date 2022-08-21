package transfer

import (
	"time"

	"github.com/pkg/errors"
)

var (
	ErrMissingSender         = errors.New("missing sender")
	ErrMissingReceiver       = errors.New("missing receiver")
	ErrSameSenderAndReceiver = errors.New("sender and receiver is the same wallet")
	ErrNonPositiveAmount     = errors.New("amount should be greater then 0")
	ErrNotEnoughMoney        = errors.New("sender does not have enough 'money' for transfer")
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
	if err := dto.validate(); err != nil {
		return nil, err
	}
	dto.Sender.Balance -= dto.Amount
	dto.Receiver.Balance += dto.Amount
	dto.Timestamp = timestamp
	return dto.toModel(), nil
}
