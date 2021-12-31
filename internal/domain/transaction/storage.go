package transaction

import "context"

type Storage interface {
	GetByID(context.Context, int64) (*TransactionDTO, error)
	GetAll(context.Context, int, int) ([]*TransactionDTO, error)
}
