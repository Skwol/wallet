package account

import (
	"context"
)

type Service interface {
	Create(context.Context, *CreateAccountDTO) (*AccountDTO, error)
	GetByID(context.Context, int64) (*AccountDTO, error)
	GetAll(ctx context.Context, limit int, offset int) ([]*AccountDTO, error)
}

type service struct {
	storage Storage
}

func NewService(storage Storage) (Service, error) {
	return &service{storage: storage}, nil
}

func (s *service) Create(context.Context, *CreateAccountDTO) (*AccountDTO, error) {
	return nil, nil
}

func (s *service) GetByID(context.Context, int64) (*AccountDTO, error) {
	return nil, nil
}

func (s *service) GetAll(ctx context.Context, limit int, offset int) ([]*AccountDTO, error) {
	return s.storage.GetAll(ctx, limit, offset)
}
