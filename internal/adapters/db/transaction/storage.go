package transaction

import (
	"context"

	"github.com/skwol/wallet/internal/domain/transaction"
	"github.com/skwol/wallet/pkg/client/pgdb"
)

type transactionStorage struct {
	db *pgdb.PGDB
}

func NewStorage(db *pgdb.PGDB) (transaction.Storage, error) {
	return &transactionStorage{db: db}, nil
}

func (as *transactionStorage) Create(ctx context.Context, acct *transaction.Transaction) (*transaction.Transaction, error) {
	return nil, nil
}

func (as *transactionStorage) GetByID(ctx context.Context, id int64) (*transaction.Transaction, error) {
	return nil, nil
}

func (as *transactionStorage) GetAll(ctx context.Context, limit int64, offset int64) ([]*transaction.Transaction, error) {
	return nil, nil
}
