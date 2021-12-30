package wallet

import (
	"context"
	"fmt"
)

type Service interface {
	Create(context.Context, *CreateWalletDTO) (*WalletDTO, error)
	GetByID(context.Context, int64) (*WalletDTO, error)
	GetAll(ctx context.Context, limit int, offset int) ([]*WalletDTO, error)
	Update(context.Context, int64, *UpdateWalletDTO) (*WalletDTO, error)
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

func (s *service) Update(ctx context.Context, id int64, walletDTO *UpdateWalletDTO) (*WalletDTO, error) {
	walletInDB, err := s.storage.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting wallet from db: %w", err)
	}
	if walletInDB == nil {
		return nil, fmt.Errorf("missing wallet db")
	}
	wallet, err := walletInDB.toModel().Update(walletDTO)
	if err != nil {
		return nil, fmt.Errorf("error updating wallet model: %w", err)
	}

	result := wallet.toDTO()
	if err := s.storage.Update(ctx, result); err != nil {
		return nil, fmt.Errorf("error updating wallet in db: %w", err)
	}
	return result, nil
}
