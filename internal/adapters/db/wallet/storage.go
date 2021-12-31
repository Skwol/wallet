package wallet

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/skwol/wallet/internal/domain/wallet"
	"github.com/skwol/wallet/pkg/client/pgdb"
)

type dbWallet struct {
	ID      int64   `json:"id,omitempty"`
	Name    string  `json:"name,omitempty"`
	Balance float64 `json:"balance,omitempty"`
}

func (db dbWallet) ToDTO() *wallet.WalletDTO {
	return &wallet.WalletDTO{
		ID:      db.ID,
		Name:    db.Name,
		Balance: db.Balance,
	}
}

type walletStorage struct {
	db *pgdb.PGDB
}

func NewStorage(db *pgdb.PGDB) (wallet.Storage, error) {
	return &walletStorage{db: db}, nil
}

func (as *walletStorage) Create(ctx context.Context, dto *wallet.WalletDTO) (*wallet.WalletDTO, error) {
	tx, err := as.db.Conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}

	row := tx.QueryRow("INSERT INTO wallet (name, balance) VALUES ($1, $2) RETURNING id;", dto.Name, dto.Balance)

	if err = row.Scan(&dto.ID); err != nil {
		tx.Rollback()
		return nil, err
	}

	for _, tran := range dto.TransactionsToApply {
		_, err = tx.ExecContext(ctx, "INSERT INTO transaction (sender_id, receiver_id, amount, date, tran_type) VALUES ($1, $1, $2, $3, $4);", dto.ID, tran.Amount, tran.Timestamp, tran.Type)
		if err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("error inserting transaction: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return nil, fmt.Errorf("error committing transction: %w", err)
	}

	return dto, nil
}

func (as *walletStorage) GetByID(ctx context.Context, id int64) (*wallet.WalletDTO, error) {
	query := `SELECT id, name, balance FROM wallet WHERE id = $1;`
	row := as.db.Conn.QueryRow(query, id)
	var walletInDB dbWallet
	switch err := row.Scan(&walletInDB.ID, &walletInDB.Name, &walletInDB.Balance); err {
	case sql.ErrNoRows:
		return nil, nil
	default:
		return walletInDB.ToDTO(), err
	}
}

func (as *walletStorage) GetByName(ctx context.Context, name string) (*wallet.WalletDTO, error) {
	query := `SELECT id FROM wallet WHERE name = $1;`
	row := as.db.Conn.QueryRow(query, name)
	var walletInDB dbWallet
	switch err := row.Scan(&walletInDB.ID); err {
	case sql.ErrNoRows:
		return nil, nil
	default:
		return walletInDB.ToDTO(), err
	}
}

func (as *walletStorage) GetAll(ctx context.Context, limit int, offset int) ([]*wallet.WalletDTO, error) {
	var list []*wallet.WalletDTO
	rows, err := as.db.Conn.Query("SELECT id, name, balance FROM wallet ORDER BY ID ASC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return list, err
	}
	var wallet dbWallet
	for rows.Next() {
		if err := rows.Scan(&wallet.ID, &wallet.Name, &wallet.Balance); err != nil {
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

	_, err = tx.ExecContext(ctx, "UPDATE wallet SET name=$1, balance=$2 WHERE id=$3;", walletDTO.Name, walletDTO.Balance, walletDTO.ID)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("error updating wallet: %w", err)
	}

	for _, tran := range walletDTO.TransactionsToApply {
		_, err = tx.ExecContext(ctx, "INSERT INTO transaction (sender_id, receiver_id, amount, date, tran_type) VALUES ($1, $1, $2, $3, $4);", walletDTO.ID, tran.Amount, tran.Timestamp, tran.Type)
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
