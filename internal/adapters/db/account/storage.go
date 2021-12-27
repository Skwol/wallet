package account

import (
	"context"

	"github.com/skwol/wallet/internal/domain/account"
)

type accountStorage struct{}

func NewStorage() account.Storage {
	return &accountStorage{}
}

func (as *accountStorage) Create(ctx context.Context, acct *account.Account) (*account.Account, error) {
	return nil, nil
}

func (as *accountStorage) GetByID(ctx context.Context, id int64) (*account.Account, error) {
	return nil, nil
}

func (as *accountStorage) GetAll(ctx context.Context, limit int64, offset int64) ([]*account.Account, error) {
	return nil, nil
}
