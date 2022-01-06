package common

import (
	"context"
)

type Service interface {
	GenerateFakeData(context.Context, int) error
}

type service struct {
	storage Storage
}

func NewService(storage Storage) (Service, error) {
	return &service{storage: storage}, nil
}

func (s *service) GenerateFakeData(ctx context.Context, numberOfRecordsToCreate int) error {
	return s.storage.GenerateFakeData(ctx, numberOfRecordsToCreate)
}
