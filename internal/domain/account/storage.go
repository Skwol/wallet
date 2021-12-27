package account

import "context"

type Storage interface {
	Create(context.Context, *Account) (*Account, error)
	GetByID(context.Context, int64) (*Account, error)
	GetAll(context.Context, int64, int64) ([]*Account, error)
}
