package transaction

import "context"

type Storage interface {
	Create(context.Context, *TransactionDTO) (*TransactionDTO, error)
	GetByID(context.Context, int64) (*TransactionDTO, error)
	GetAll(context.Context, int, int) ([]*TransactionDTO, error)
}
