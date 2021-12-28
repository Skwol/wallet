package account

import (
	"context"

	"github.com/skwol/wallet/internal/domain/account"
	"github.com/skwol/wallet/pkg/client/pgdb"
)

type accountStorage struct {
	db *pgdb.PGDB
}

func NewStorage(db *pgdb.PGDB) (account.Storage, error) {
	return &accountStorage{db: db}, nil
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
