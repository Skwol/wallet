package account

import (
	"context"

	"github.com/skwol/wallet/internal/domain/account"
	"github.com/skwol/wallet/pkg/client/pgdb"
)

type dbAccount struct {
	ID       int64  `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
}

func (db dbAccount) ToDTO() *account.AccountDTO {
	return &account.AccountDTO{
		ID:       db.ID,
		Username: db.Username,
	}
}

type accountStorage struct {
	db *pgdb.PGDB
}

func NewStorage(db *pgdb.PGDB) (account.Storage, error) {
	return &accountStorage{db: db}, nil
}

func (as *accountStorage) Create(ctx context.Context, acct *account.AccountDTO) (*account.AccountDTO, error) {
	return nil, nil
}

func (as *accountStorage) GetByID(ctx context.Context, id int64) (*account.AccountDTO, error) {
	return nil, nil
}

func (as *accountStorage) GetAll(ctx context.Context, limit int, offset int) ([]*account.AccountDTO, error) {
	var list []*account.AccountDTO
	rows, err := as.db.Conn.Query("SELECT * FROM account ORDER BY ID ASC LIMIT $1 OFFSET $2;", limit, offset)
	if err != nil {
		return nil, err
	}
	var account dbAccount
	for rows.Next() {
		if err := rows.Scan(&account.ID, &account.Username); err != nil {
			return nil, err
		}
		list = append(list, account.ToDTO())
	}
	return list, nil
}
