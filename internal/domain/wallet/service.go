package wallet

import (
	"context"
)

type Service interface {
	Create(context.Context, *CreateWalletDTO) (*Wallet, error)
	GetByID(context.Context, int64) (*Wallet, error)
	GetAll(ctx context.Context, limit int, offset int) ([]*Wallet, error)
	Update(context.Context, *UpdateWalletDTO) (*Wallet, error)
}

type service struct {
	storage Storage
}

func NewService(storage Storage) (Service, error) {
	return &service{storage: storage}, nil
}

func (s *service) Create(context.Context, *CreateWalletDTO) (*Wallet, error) {
	return nil, nil
}

func (s *service) GetByID(context.Context, int64) (*Wallet, error) {
	return nil, nil
}

func (s *service) GetAll(ctx context.Context, limit int, offset int) ([]*Wallet, error) {
	return nil, nil
}

func (s *service) Update(context.Context, *UpdateWalletDTO) (*Wallet, error) {
	return nil, nil
}
