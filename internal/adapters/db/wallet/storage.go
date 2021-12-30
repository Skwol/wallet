package wallet

import (
	"context"
	"database/sql"
	"fmt"

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
	wallet := wallet.WalletDTO{}
	query := `SELECT * FROM wallet WHERE id = $1;`
	row := as.db.Conn.QueryRow(query, id)
	switch err := row.Scan(&wallet.ID, &wallet.Name, &wallet.AccountID, &wallet.Balance); err {
	case sql.ErrNoRows:
		return nil, nil
	default:
		return &wallet, err
	}
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

func (as *walletStorage) Update(ctx context.Context, walletDTO *wallet.WalletDTO) error {
	tx, err := as.db.Conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction: %w", err)
	}

	var balance float64
	row := tx.QueryRowContext(ctx, "SELECT balance FROM wallet WHERE id = $1 FOR UPDATE;", walletDTO.ID)
	err = row.Scan(&balance)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error getting balance for wallet: %w", err)
	}

	_, err = tx.ExecContext(ctx, "UPDATE wallet SET name=$1, balance=$2 WHERE id=$3;", walletDTO.Name, walletDTO.Balance, walletDTO.ID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error updating wallet: %w", err)
	}

	for _, tran := range walletDTO.TransactionsToApply {
		_, err = tx.ExecContext(ctx, "INSERT INTO transaction (sender_id, receiver_id, amount, date, tran_type) VALUES ($1, $1, $2, current_timestamp, $3);", walletDTO.ID, tran.Amount, tran.Type)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("error inserting transaction: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error committing transction: %w", err)
	}

	return nil
}
