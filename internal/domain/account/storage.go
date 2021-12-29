package account

import "context"

type Storage interface {
	Create(context.Context, *AccountDTO) (*AccountDTO, error)
	GetByID(context.Context, int64) (*AccountDTO, error)
	GetAll(context.Context, int, int) ([]*AccountDTO, error)
}
