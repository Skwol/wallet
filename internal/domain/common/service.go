package common

import (
	"context"

	"github.com/skwol/wallet/pkg/logging"
)

type Service interface {
	GenerateFakeData(context.Context, int) error
}

type service struct {
	storage Storage
	logger  logging.Logger
}

func NewService(storage Storage, logger logging.Logger) (Service, error) {
	return &service{storage: storage, logger: logger}, nil
}

func (s *service) GenerateFakeData(ctx context.Context, numberOfRecordsToCreate int) error {
	return s.storage.GenerateFakeData(ctx, numberOfRecordsToCreate)
}
