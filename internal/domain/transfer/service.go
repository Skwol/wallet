package transfer

import (
	"context"
	"fmt"

	"github.com/skwol/wallet/pkg/logging"
)

type Service interface {
	Create(context.Context, *CreateTransferDTO) (TransferDTO, error)
}

type service struct {
	storage Storage
}

func NewService(storage Storage) (Service, error) {
	return &service{storage: storage}, nil
}

func (s *service) Create(ctx context.Context, dto *CreateTransferDTO) (TransferDTO, error) {
	logger := logging.GetLogger()

	walletSender, err := s.storage.GetWallet(ctx, dto.Sender.ID)
	if err != nil {
		logger.Errorf("error getting sender wallet from db: %s", err.Error())
		return TransferDTO{}, fmt.Errorf("error getting sender wallet from db: %w", err)
	}
	if walletSender.ID == 0 {
		logger.Errorf("missing sender wallet in db")
		return TransferDTO{}, fmt.Errorf("missing sender wallet in db")
	}

	walletReceiver, err := s.storage.GetWallet(ctx, dto.Receiver.ID)
	if err != nil {
		logger.Errorf("error getting receiver wallet from db: %s", err.Error())
		return TransferDTO{}, fmt.Errorf("error getting receiver wallet from db: %w", err)
	}
	if walletReceiver.ID == 0 {
		logger.Errorf("missing receiver wallet in db")
		return TransferDTO{}, fmt.Errorf("missing receiver wallet in db")
	}
	dto.Sender = walletSender
	dto.Receiver = walletReceiver

	transferModel, err := createTransfer(dto)
	if err != nil {
		logger.Errorf("error creating transfer model: %s", err.Error())
		return TransferDTO{}, fmt.Errorf("error creating transfer model: %w", err)
	}
	if transferModel == nil {
		logger.Errorf("transfer model was not created")
		return TransferDTO{}, fmt.Errorf("transfer model was not created")
	}
	result, err := s.storage.Create(ctx, &transferModel.toDTO().CreateTransferDTO)
	if err != nil {
		logger.Errorf("error creating transfer in db: %s", err.Error())
		return TransferDTO{}, fmt.Errorf("error creating wallet in db: %w", err)
	}
	if result.ID == 0 {
		logger.Errorf("empty transfer returned from db")
		return TransferDTO{}, fmt.Errorf("empty transfer returned from db")
	}
	return result, nil
}
