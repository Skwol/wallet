package wallet

import (
	"context"

	"github.com/skwol/wallet/internal/domain/wallet"
	"github.com/skwol/wallet/pkg/client/pgdb"
)

type dbWallet struct {
	ID        int64   `json:"id,omitempty"`
	Name      string  `json:"name,omitempty"`
	AccountID int64   `json:"account_id,omitempty"`
	Balance   float64 `json:"balance,omitempty"`
}

func (db dbWallet) ToDTO() *wallet.WalletDTO {
	return &wallet.WalletDTO{
		ID:        db.ID,
		Name:      db.Name,
		AccountID: db.AccountID,
		Balance:   db.Balance,
	}
}

type walletStorage struct {
	db *pgdb.PGDB
}

func NewStorage(db *pgdb.PGDB) (wallet.Storage, error) {
	return &walletStorage{db: db}, nil
}

func (as *walletStorage) Create(ctx context.Context, acct *wallet.WalletDTO) (*wallet.WalletDTO, error) {
	return nil, nil
}

func (as *walletStorage) GetByID(ctx context.Context, id int64) (*wallet.WalletDTO, error) {
	return nil, nil
}

func (as *walletStorage) GetAll(ctx context.Context, limit int, offset int) ([]*wallet.WalletDTO, error) {
	var list []*wallet.WalletDTO
	rows, err := as.db.Conn.Query("SELECT * FROM wallet ORDER BY ID ASC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return list, err
	}
	var wallet dbWallet
	for rows.Next() {
		if err := rows.Scan(&wallet.ID, &wallet.Name, &wallet.AccountID, &wallet.Balance); err != nil {
			return nil, err
		}
		list = append(list, wallet.ToDTO())
	}
	return list, nil
}

func (as *walletStorage) Update(context.Context, *wallet.WalletDTO) (*wallet.WalletDTO, error) {
	return nil, nil
}
