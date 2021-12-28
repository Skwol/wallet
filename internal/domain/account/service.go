package account

import (
	"context"
)

type Service interface {
	Create(context.Context, *CreateAccountDTO) (*Account, error)
	GetByID(context.Context, int64) (*Account, error)
	GetAll(ctx context.Context, limit int, offset int) ([]*Account, error)
}

type service struct {
	storage Storage
}

func NewService(storage Storage) (Service, error) {
	return &service{storage: storage}, nil
}

func (s *service) Create(context.Context, *CreateAccountDTO) (*Account, error) {
	return nil, nil
}

func (s *service) GetByID(context.Context, int64) (*Account, error) {
	return nil, nil
}

func (s *service) GetAll(ctx context.Context, limit int, offset int) ([]*Account, error) {
	return nil, nil
}
