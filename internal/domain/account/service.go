package account

import "context"

type Service interface {
	Create(context.Context, *Account) (*Account, error)
	GetByID(context.Context, int64) (*Account, error)
	GetAll(ctx context.Context, limit int, offset int) ([]*Account, error)
}

type service struct {
	storage Storage
}

func NewService(storage Storage) Service {
	return &service{storage: storage}
}

func (s *service) Create(context.Context, *Account) (*Account, error) {
	return nil, nil
}

func (s *service) GetByID(context.Context, int64) (*Account, error) {
	return nil, nil
}

func (s *service) GetAll(ctx context.Context, limit int, offset int) ([]*Account, error) {
	return nil, nil
}
