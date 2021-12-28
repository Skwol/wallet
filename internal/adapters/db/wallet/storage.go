package wallet

import (
	"context"

	"github.com/skwol/wallet/internal/domain/wallet"
	"github.com/skwol/wallet/pkg/client/pgdb"
)

type walletStorage struct {
	db *pgdb.PGDB
}

func NewStorage(db *pgdb.PGDB) (wallet.Storage, error) {
	return &walletStorage{db: db}, nil
}

func (as *walletStorage) Create(ctx context.Context, acct *wallet.Wallet) (*wallet.Wallet, error) {
	return nil, nil
}

func (as *walletStorage) GetByID(ctx context.Context, id int64) (*wallet.Wallet, error) {
	return nil, nil
}

func (as *walletStorage) GetAll(ctx context.Context, limit int64, offset int64) ([]*wallet.Wallet, error) {
	return nil, nil
}

func (as *walletStorage) Update(context.Context, *wallet.Wallet) (*wallet.Wallet, error) {
	return nil, nil
}
