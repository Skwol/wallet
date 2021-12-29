package wallet

import (
	"context"
)

type Service interface {
	Create(context.Context, *CreateWalletDTO) (*WalletDTO, error)
	GetByID(context.Context, int64) (*WalletDTO, error)
	GetAll(ctx context.Context, limit int, offset int) ([]*WalletDTO, error)
	Update(context.Context, *UpdateWalletDTO) (*WalletDTO, error)
}

type service struct {
	storage Storage
}

func NewService(storage Storage) (Service, error) {
	return &service{storage: storage}, nil
}

func (s *service) Create(context.Context, *CreateWalletDTO) (*WalletDTO, error) {
	return nil, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (*WalletDTO, error) {
	return s.storage.GetByID(ctx, id)
}

func (s *service) GetAll(ctx context.Context, limit int, offset int) ([]*WalletDTO, error) {
	return s.storage.GetAll(ctx, limit, offset)
}

func (s *service) Update(context.Context, *UpdateWalletDTO) (*WalletDTO, error) {
	return nil, nil
}
