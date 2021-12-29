package transaction

import (
	"context"
)

type Service interface {
	Create(context.Context, *CreateTransactionDTO) (*TransactionDTO, error)
	GetByID(context.Context, int64) (*TransactionDTO, error)
	GetAll(ctx context.Context, limit int, offset int) ([]*TransactionDTO, error)
}

type service struct {
	storage Storage
}

func NewService(storage Storage) (Service, error) {
	return &service{storage: storage}, nil
}

func (s *service) Create(context.Context, *CreateTransactionDTO) (*TransactionDTO, error) {
	return nil, nil
}

func (s *service) GetByID(context.Context, int64) (*TransactionDTO, error) {
	return nil, nil
}

func (s *service) GetAll(ctx context.Context, limit int, offset int) ([]*TransactionDTO, error) {
	return s.storage.GetAll(ctx, limit, offset)
}
