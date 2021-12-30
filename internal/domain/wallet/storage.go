package wallet

import "context"

type Storage interface {
	Create(context.Context, *WalletDTO) (*WalletDTO, error)
	GetByID(context.Context, int64) (*WalletDTO, error)
	GetAll(context.Context, int, int) ([]*WalletDTO, error)
	Update(context.Context, *WalletDTO) error
}
