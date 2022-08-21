package transfer

import (
	"github.com/skwol/wallet/internal/domain/transfer"
)

func newTransfer(dto transfer.DTO) Transfer {
	return Transfer{
		Id:        int(dto.ID),
		Amount:    float32(dto.Amount),
		Timestamp: &dto.Timestamp,
		Sender:    newWallet(dto.Sender),
		Receiver:  newWallet(dto.Receiver),
	}
}

func (w CreateTransferRequest) toCreateRequest() transfer.CreateTransferDTO {
	return transfer.CreateTransferDTO{
		Amount:   float64(w.Amount),
		Sender:   transfer.WalletDTO{ID: int64(w.SenderId)},
		Receiver: transfer.WalletDTO{ID: int64(w.ReceiverId)},
	}
}

func newWallet(dto transfer.WalletDTO) Wallet {
	return Wallet{
		Id:      int(dto.ID),
		Balance: float32(dto.Balance),
	}
}

func (w Wallet) toRequest() transfer.WalletDTO {
	return transfer.WalletDTO{
		ID: int64(w.Id),
	}
}
