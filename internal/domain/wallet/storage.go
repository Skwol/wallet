package wallet

import "context"

type Storage interface {
	Create(context.Context, *Wallet) (*Wallet, error)
	GetByID(context.Context, int64) (*Wallet, error)
	GetAll(context.Context, int64, int64) ([]*Wallet, error)
	Update(context.Context, *Wallet) (*Wallet, error)
}
