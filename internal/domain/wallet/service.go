package wallet

import (
	"context"

	"github.com/pkg/errors"

	"github.com/skwol/wallet/pkg/clock"
	"github.com/skwol/wallet/pkg/logging"
)

type Service interface {
	Create(context.Context, *CreateWalletDTO) (DTO, error)
	GetByID(context.Context, int64) (DTO, error)
	GetByIDWithTransactions(context.Context, int64, int, int) (DTO, error)
	GetAll(ctx context.Context, limit int, offset int) ([]DTO, error)
	Update(context.Context, int64, *UpdateWalletDTO) (DTO, error)
}

type service struct {
	storage Storage
	logger  logging.Logger
	clk     clock.Clock
}

func NewService(storage Storage, logger logging.Logger, clk clock.Clock) (Service, error) {
	return &service{storage: storage, logger: logger, clk: clk}, nil
}

func (s *service) Create(ctx context.Context, dto *CreateWalletDTO) (DTO, error) {
	var result DTO
	dbWallet, err := s.storage.GetByName(ctx, dto.Name)
	if err != nil {
		s.logger.Errorf("error getting wallet from db: %s", err.Error())
		return result, errors.Wrap(err, "error getting wallet from db")
	}
	if dbWallet.ID != 0 {
		s.logger.Errorf("wallet with name %s already exist", dto.Name)
		return result, errors.Errorf("wallet with name %s already exist", dto.Name)
	}
	walletModel, err := newWallet(dto, s.clk.Now())
	if err != nil {
		s.logger.Errorf("error creating wallet model: %s", err.Error())
		return result, errors.Wrap(err, "error creating wallet model")
	}
	if walletModel == nil {
		s.logger.Errorf("wallet model was not created")
		return result, errors.Wrap(err, "wallet model was not created")
	}
	result, err = s.storage.Create(ctx, walletModel.toDTO())
	if err != nil {
		s.logger.Errorf("error creating wallet in db: %s", err.Error())
		return result, errors.Wrap(err, "error creating wallet in db")
	}
	if result.ID == 0 {
		s.logger.Errorf("empty wallet returned from db")
		return result, errors.Wrap(err, "empty wallet returned from db")
	}
	return result, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (DTO, error) {
	return s.storage.GetByID(ctx, id)
}

func (s *service) GetByIDWithTransactions(ctx context.Context, id int64, limit, offset int) (DTO, error) {
	return s.storage.GetByIDWithTransactions(ctx, id, limit, offset)
}

func (s *service) GetAll(ctx context.Context, limit int, offset int) ([]DTO, error) {
	return s.storage.GetAll(ctx, limit, offset)
}

func (s *service) Update(ctx context.Context, id int64, walletDTO *UpdateWalletDTO) (DTO, error) {
	var result DTO

	walletInDB, err := s.storage.GetByID(ctx, id)
	if err != nil {
		s.logger.Errorf("error getting wallet from db: %s", err.Error())
		return result, errors.Wrap(err, "error getting wallet from db")
	}
	if walletInDB.ID == 0 {
		s.logger.Errorf("missing wallet in db")
		return result, errors.New("missing wallet db")
	}
	walletModel := walletInDB.toModel()
	wallet, err := walletModel.Update(walletDTO, s.clk.Now())
	if err != nil {
		s.logger.Errorf("error updating wallet model: %s", err.Error())
		return result, errors.Wrap(err, "error updating wallet model")
	}

	result = wallet.toDTO()
	if err := s.storage.Update(ctx, result); err != nil {
		s.logger.Errorf("error updating wallet in db: %s", err.Error())
		return result, errors.Wrap(err, "error updating wallet in db")
	}
	return result, nil
}
