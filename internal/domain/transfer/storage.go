package transfer

import "context"

type Storage interface {
	Create(context.Context, *CreateTransferDTO) (DTO, error)
	GetWallet(context.Context, int64) (WalletDTO, error)
}
