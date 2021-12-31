package transaction

import (
	"context"
)

type Service interface {
	GetByID(context.Context, int64) (*TransactionDTO, error)
	GetAll(ctx context.Context, limit int, offset int) ([]*TransactionDTO, error)
}

type service struct {
	storage Storage
}

func NewService(storage Storage) (Service, error) {
	return &service{storage: storage}, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (*TransactionDTO, error) {
	return s.storage.GetByID(ctx, id)
}

func (s *service) GetAll(ctx context.Context, limit int, offset int) ([]*TransactionDTO, error) {
	return s.storage.GetAll(ctx, limit, offset)
}
