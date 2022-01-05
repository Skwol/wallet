package wallet

import (
	"context"
	"fmt"

	"github.com/skwol/wallet/pkg/logging"
)

type Service interface {
	Create(context.Context, CreateWalletDTO) (WalletDTO, error)
	GetByID(context.Context, int64) (WalletDTO, error)
	GetByIDWithTransactions(context.Context, int64, int, int) (WalletDTO, error)
	GetAll(ctx context.Context, limit int, offset int) ([]WalletDTO, error)
	Update(context.Context, int64, *UpdateWalletDTO) (WalletDTO, error)
}

type service struct {
	storage Storage
}

func NewService(storage Storage) (Service, error) {
	return &service{storage: storage}, nil
}

func (s *service) Create(ctx context.Context, dto CreateWalletDTO) (WalletDTO, error) {
	var result WalletDTO
	logger := logging.GetLogger()
	dbWallet, err := s.storage.GetByName(ctx, dto.Name)
	if err != nil {
		logger.Errorf("error getting wallet from db: %s", err.Error())
		return result, fmt.Errorf("error getting wallet from db: %w", err)
	}
	if dbWallet.ID != 0 {
		logger.Errorf("wallet with name %s already exist", dto.Name)
		return result, fmt.Errorf("wallet with name %s already exist", dto.Name)
	}
	walletModel, err := newWallet(&dto)
	if err != nil {
		logger.Errorf("error creating wallet model: %s", err.Error())
		return result, fmt.Errorf("error creating wallet model: %w", err)
	}
	if walletModel == nil {
		logger.Errorf("wallet model was not created")
		return result, fmt.Errorf("wallet model was not created")
	}
	result, err = s.storage.Create(ctx, walletModel.toDTO())
	if err != nil {
		logger.Errorf("error creating wallet in db: %s", err.Error())
		return result, fmt.Errorf("error creating wallet in db: %w", err)
	}
	if result.ID == 0 {
		logger.Errorf("empty wallet returned from db")
		return result, fmt.Errorf("empty wallet returned from db")
	}
	return result, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (WalletDTO, error) {
	return s.storage.GetByID(ctx, id)
}

func (s *service) GetByIDWithTransactions(ctx context.Context, id int64, limit, offset int) (WalletDTO, error) {
	return s.storage.GetByIDWithTransactions(ctx, id, limit, offset)
}

func (s *service) GetAll(ctx context.Context, limit int, offset int) ([]WalletDTO, error) {
	return s.storage.GetAll(ctx, limit, offset)
}

func (s *service) Update(ctx context.Context, id int64, walletDTO *UpdateWalletDTO) (WalletDTO, error) {
	var result WalletDTO
	logger := logging.GetLogger()

	walletInDB, err := s.storage.GetByID(ctx, id)
	if err != nil {
		logger.Errorf("error getting wallet from db: %s", err.Error())
		return result, fmt.Errorf("error getting wallet from db: %w", err)
	}
	if walletInDB.ID == 0 {
		logger.Errorf("missing wallet in db")
		return result, fmt.Errorf("missing wallet db")
	}
	wallet, err := walletInDB.toModel().Update(walletDTO)
	if err != nil {
		logger.Errorf("error updating wallet model: %s", err.Error())
		return result, fmt.Errorf("error updating wallet model: %w", err)
	}

	result = wallet.toDTO()
	if err := s.storage.Update(ctx, result); err != nil {
		logger.Errorf("error updating wallet in db: %s", err.Error())
		return result, fmt.Errorf("error updating wallet in db: %w", err)
	}
	return result, nil
}
