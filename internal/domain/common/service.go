package common

import (
	"context"
)

type Service interface {
	GenerateFakeData(context.Context) error
}

type service struct {
	storage Storage
}

func NewService(storage Storage) (Service, error) {
	return &service{storage: storage}, nil
}

func (s *service) GenerateFakeData(ctx context.Context) error {
	return s.storage.GenerateFakeData(ctx)
}
