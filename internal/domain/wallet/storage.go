package wallet

import "context"

type Storage interface {
	Create(context.Context, DTO) (DTO, error)
	GetByID(context.Context, int64) (DTO, error)
	GetByIDWithTransactions(context.Context, int64, int, int) (DTO, error)
	GetByName(context.Context, string) (DTO, error)
	GetAll(context.Context, int, int) ([]DTO, error)
	Update(context.Context, DTO) error
}
