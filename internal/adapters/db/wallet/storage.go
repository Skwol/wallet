package wallet

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/skwol/wallet/internal/domain/wallet"
	"github.com/skwol/wallet/pkg/client/pgdb"
)

type dbWallet struct {
	ID      int64   `json:"id,omitempty"`
	Name    string  `json:"name,omitempty"`
	Balance float64 `json:"balance,omitempty"`
}

func (db dbWallet) ToDTO() wallet.WalletDTO {
	return wallet.WalletDTO{
		ID:      db.ID,
		Name:    db.Name,
		Balance: db.Balance,
	}
}

type dbTransaction struct {
	ID         int64           `json:"id,omitempty"`
	SenderID   int64           `json:"sender_id,omitempty"`
	ReceiverID int64           `json:"receiver_id,omitempty"`
	Amount     float64         `json:"amount,omitempty"`
	Timestamp  time.Time       `json:"timestamp,omitempty"`
	Type       wallet.TranType `json:"type,omitempty"`
}

func (db dbTransaction) ToDTO() *wallet.TransactionDTO {
	return &wallet.TransactionDTO{
		ID:         db.ID,
		SenderID:   db.SenderID,
		ReceiverID: db.ReceiverID,
		Amount:     db.Amount,
		Timestamp:  db.Timestamp,
		Type:       db.Type,
	}
}

type walletStorage struct {
	db *pgdb.PGDB
}

func NewStorage(db *pgdb.PGDB) (wallet.Storage, error) {
	return &walletStorage{db: db}, nil
}

func (as *walletStorage) Create(ctx context.Context, dto wallet.WalletDTO) (wallet.WalletDTO, error) {
	tx, err := as.db.Conn.BeginTx(ctx, nil)
	if err != nil {
		return dto, err
	}

	row := tx.QueryRow("INSERT INTO wallet (name, balance) VALUES ($1, $2) RETURNING id;", dto.Name, dto.Balance)

	if err = row.Scan(&dto.ID); err != nil {
		tx.Rollback()
		return dto, err
	}

	for _, tran := range dto.TransactionsToApply {
		_, err = tx.ExecContext(ctx, "INSERT INTO transaction (sender_id, receiver_id, amount, date, tran_type) VALUES ($1, $1, $2, $3, $4);", dto.ID, tran.Amount, tran.Timestamp, tran.Type)
		if err != nil {
			tx.Rollback()
			return dto, fmt.Errorf("error inserting transaction: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return dto, fmt.Errorf("error committing transction: %w", err)
	}

	return dto, nil
}

func (as *walletStorage) GetByID(ctx context.Context, id int64) (wallet.WalletDTO, error) {
	query := `SELECT id, name, balance FROM wallet WHERE id = $1;`
	row := as.db.Conn.QueryRow(query, id)
	var walletInDB dbWallet
	switch err := row.Scan(&walletInDB.ID, &walletInDB.Name, &walletInDB.Balance); err {
	case sql.ErrNoRows:
		return wallet.WalletDTO{}, nil
	default:
		return walletInDB.ToDTO(), err
	}
}

func (as *walletStorage) GetByIDWithTransactions(ctx context.Context, id int64, limit int, offset int) (wallet.WalletDTO, error) {
	tx, err := as.db.Conn.BeginTx(ctx, nil)
	if err != nil {
		return wallet.WalletDTO{}, fmt.Errorf("error beginning transaction: %w", err)
	}
	query := `SELECT id, name, balance FROM wallet WHERE id = $1;`
	row := as.db.Conn.QueryRow(query, id)
	var walletInDB dbWallet
	if err := row.Scan(&walletInDB.ID, &walletInDB.Name, &walletInDB.Balance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return wallet.WalletDTO{}, nil
		}
		return wallet.WalletDTO{}, err
	}

	query = "SELECT id, sender_id, receiver_id, amount, date, tran_type FROM transaction WHERE sender_id = $1 OR receiver_id = $1 ORDER BY ID ASC LIMIT $2 OFFSET $3"
	rows, err := as.db.Conn.Query(query, walletInDB.ID, limit, offset)
	if err != nil {
		return wallet.WalletDTO{}, err
	}
	var (
		list []*wallet.TransactionDTO
		tran dbTransaction
	)
	for rows.Next() {
		if err := rows.Scan(&tran.ID, &tran.SenderID, &tran.ReceiverID, &tran.Amount, &tran.Timestamp, &tran.Type); err != nil {
			return wallet.WalletDTO{}, err
		}
		list = append(list, tran.ToDTO())
	}

	if err := tx.Commit(); err != nil {
		return wallet.WalletDTO{}, fmt.Errorf("error committing transction: %w", err)
	}
	walletDTO := walletInDB.ToDTO()
	walletDTO.Transactions = list
	return walletDTO, nil
}

func (as *walletStorage) GetByName(ctx context.Context, name string) (wallet.WalletDTO, error) {
	query := `SELECT id FROM wallet WHERE name = $1;`
	row := as.db.Conn.QueryRow(query, name)
	var walletInDB dbWallet
	switch err := row.Scan(&walletInDB.ID); err {
	case sql.ErrNoRows:
		return wallet.WalletDTO{}, nil
	default:
		return walletInDB.ToDTO(), err
	}
}

func (as *walletStorage) GetAll(ctx context.Context, limit int, offset int) ([]wallet.WalletDTO, error) {
	var list []wallet.WalletDTO
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

func (as *walletStorage) Update(ctx context.Context, walletDTO wallet.WalletDTO) error {
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
