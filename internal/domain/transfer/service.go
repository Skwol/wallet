package transfer

import (
	"context"
	"fmt"

	"github.com/skwol/wallet/pkg/clock"
	"github.com/skwol/wallet/pkg/logging"
)

type Service interface {
	Create(context.Context, *CreateTransferDTO) (DTO, error)
}

type service struct {
	storage Storage
	logger  logging.Logger
	clk     clock.Clock
}

func NewService(storage Storage, logger logging.Logger, clk clock.Clock) (Service, error) {
	return &service{storage: storage, logger: logger, clk: clk}, nil
}

func (s *service) Create(ctx context.Context, dto *CreateTransferDTO) (DTO, error) {
	walletSender, err := s.storage.GetWallet(ctx, dto.Sender.ID)
	if err != nil {
		s.logger.Errorf("error getting sender wallet from db: %s", err.Error())
		return DTO{}, fmt.Errorf("error getting sender wallet from db: %w", err)
	}
	if walletSender.ID == 0 {
		s.logger.Errorf("missing sender wallet in db")
		return DTO{}, fmt.Errorf("missing sender wallet in db")
	}

	walletReceiver, err := s.storage.GetWallet(ctx, dto.Receiver.ID)
	if err != nil {
		s.logger.Errorf("error getting receiver wallet from db: %s", err.Error())
		return DTO{}, fmt.Errorf("error getting receiver wallet from db: %w", err)
	}
	if walletReceiver.ID == 0 {
		s.logger.Errorf("missing receiver wallet in db")
		return DTO{}, fmt.Errorf("missing receiver wallet in db")
	}
	dto.Sender = walletSender
	dto.Receiver = walletReceiver

	transferModel, err := createTransfer(dto, s.clk.Now())
	if err != nil {
		s.logger.Errorf("error creating transfer model: %s", err.Error())
		return DTO{}, fmt.Errorf("error creating transfer model: %w", err)
	}
	if transferModel == nil {
		s.logger.Errorf("transfer model was not created")
		return DTO{}, fmt.Errorf("transfer model was not created")
	}
	result, err := s.storage.Create(ctx, &transferModel.toDTO().CreateTransferDTO)
	if err != nil {
		s.logger.Errorf("error creating transfer in db: %s", err.Error())
		return DTO{}, fmt.Errorf("error creating wallet in db: %w", err)
	}
	if result.ID == 0 {
		s.logger.Errorf("empty transfer returned from db")
		return DTO{}, fmt.Errorf("empty transfer returned from db")
	}
	return result, nil
}
