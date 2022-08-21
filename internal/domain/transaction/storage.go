package transaction

import "context"

type Storage interface {
	GetByID(context.Context, int64) (DTO, error)
	GetAll(context.Context, int, int) ([]DTO, error)
	GetFiltered(context.Context, *FilterTransactionsDTO, int, int) ([]DTO, error)
}
