package transaction

import (
	"context"

	"github.com/skwol/wallet/pkg/logging"
)

type Service interface {
	GetByID(context.Context, int64) (DTO, error)
	GetAll(ctx context.Context, limit int, offset int) ([]DTO, error)
	GetFiltered(ctx context.Context, filter *FilterTransactionsDTO, limit int, offset int) ([]DTO, error)
}

type service struct {
	storage Storage
	logger  logging.Logger
}

func NewService(storage Storage, logger logging.Logger) (Service, error) {
	return &service{storage: storage, logger: logger}, nil
}

func (s *service) GetByID(ctx context.Context, id int64) (DTO, error) {
	return s.storage.GetByID(ctx, id)
}

func (s *service) GetAll(ctx context.Context, limit int, offset int) ([]DTO, error) {
	return s.storage.GetAll(ctx, limit, offset)
}

func (s *service) GetFiltered(ctx context.Context, filter *FilterTransactionsDTO, limit int, offset int) ([]DTO, error) {
	return s.storage.GetFiltered(ctx, filter, limit, offset)
}
