package transaction

import (
	"context"
)

type Service interface {
	Create(context.Context, *CreateTransactionDTO) (*Transaction, error)
	GetByID(context.Context, int64) (*Transaction, error)
	GetAll(ctx context.Context, limit int, offset int) ([]*Transaction, error)
}

type service struct {
	storage Storage
}

func NewService(storage Storage) (Service, error) {
	return &service{storage: storage}, nil
}

func (s *service) Create(context.Context, *CreateTransactionDTO) (*Transaction, error) {
	return nil, nil
}

func (s *service) GetByID(context.Context, int64) (*Transaction, error) {
	return nil, nil
}

func (s *service) GetAll(ctx context.Context, limit int, offset int) ([]*Transaction, error) {
	return nil, nil
}
