package transfer

import (
	"time"

	"github.com/skwol/wallet/internal/domain/transfer"
)

func newTransfer(dto transfer.DTO) Transfer {
	return Transfer{
		ID:        dto.ID,
		Amount:    dto.Amount,
		Timestamp: dto.Timestamp,
		Sender:    newWallet(dto.Sender),
		Receiver:  newWallet(dto.Receiver),
	}
}

type Transfer struct {
	ID        int64     `json:"id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
	Sender    Wallet    `json:"sender"`
	Receiver  Wallet    `json:"receiver"`
}

func (w Transfer) toCreateRequest() transfer.CreateTransferDTO {
	return transfer.CreateTransferDTO{
		Amount:    w.Amount,
		Timestamp: w.Timestamp,
		Sender:    w.Sender.toRequest(),
		Receiver:  w.Receiver.toRequest(),
	}
}

func newWallet(dto transfer.WalletDTO) Wallet {
	return Wallet{
		ID:      dto.ID,
		Balance: dto.Balance,
	}
}

type Wallet struct {
	ID      int64   `json:"id"`
	Balance float64 `json:"balance"`
}

func (w Wallet) toRequest() transfer.WalletDTO {
	return transfer.WalletDTO{
		ID: w.ID,
	}
}
