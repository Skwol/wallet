package wallet

import (
	"context"

	"github.com/skwol/wallet/internal/domain/wallet"
)

type walletStorage struct{}

func NewStorage() wallet.Storage {
	return &walletStorage{}
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
