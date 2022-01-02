package wallet

import "context"

type Storage interface {
	Create(context.Context, WalletDTO) (WalletDTO, error)
	GetByID(context.Context, int64) (WalletDTO, error)
	GetByIDWithTransactions(context.Context, int64, int, int) (WalletDTO, error)
	GetByName(context.Context, string) (WalletDTO, error)
	GetAll(context.Context, int, int) ([]WalletDTO, error)
	Update(context.Context, WalletDTO) error
}
