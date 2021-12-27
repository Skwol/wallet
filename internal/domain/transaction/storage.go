package transaction

import "context"

type Storage interface {
	Create(context.Context, *Transaction) (*Transaction, error)
	GetByID(context.Context, int64) (*Transaction, error)
	GetAll(context.Context, int64, int64) ([]*Transaction, error)
}
