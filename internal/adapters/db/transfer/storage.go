package transfer

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/skwol/wallet/internal/domain/transfer"
	"github.com/skwol/wallet/pkg/client/pgdb"
)

type dbWallet struct {
	ID      int64
	Balance float64
}

func (db dbWallet) ToDTO() transfer.WalletDTO {
	return transfer.WalletDTO{
		ID:      db.ID,
		Balance: db.Balance,
	}
}

type transferStorage struct {
	db *pgdb.PGDB
}

func NewStorage(db *pgdb.PGDB) (transfer.Storage, error) {
	return &transferStorage{db: db}, nil
}

func (ts transferStorage) Create(ctx context.Context, dto *transfer.CreateTransferDTO) (transfer.TransferDTO, error) {
	var result transfer.TransferDTO
	result.CreateTransferDTO = *dto
	tx, err := ts.db.Conn.BeginTx(ctx, nil)
	if err != nil {
		tx.Rollback()
		return result, fmt.Errorf("error beginning transaction: %w", err)
	}
	if _, err = tx.ExecContext(ctx, "UPDATE wallet SET balance=$1 WHERE id=$2;", dto.Sender.Balance, dto.Sender.ID); err != nil {
		tx.Rollback()
		return result, fmt.Errorf("error updating sender wallet: %w", err)
	}
	if _, err = tx.ExecContext(ctx, "UPDATE wallet SET balance=$1 WHERE id=$2;", dto.Receiver.Balance, dto.Receiver.ID); err != nil {
		tx.Rollback()
		return result, fmt.Errorf("error updating receiver wallet: %w", err)
	}
	row := tx.QueryRow("INSERT INTO transaction (sender_id, receiver_id, amount, date, tran_type) VALUES ($1, $2, $3, $4, 'transfer') RETURNING id;", dto.Sender.ID, dto.Receiver.ID, dto.Amount, dto.Timestamp)

	if err = row.Scan(&result.ID); err != nil {
		tx.Rollback()
		return result, err
	}
	if err = tx.Commit(); err != nil {
		return result, fmt.Errorf("error during commit")
	}
	return result, nil
}

func (ts transferStorage) GetWallet(ctx context.Context, id int64) (transfer.WalletDTO, error) {
	query := `SELECT id, balance FROM wallet WHERE id = $1;`
	row := ts.db.Conn.QueryRow(query, id)
	var walletInDB dbWallet
	switch err := row.Scan(&walletInDB.ID, &walletInDB.Balance); err {
	case sql.ErrNoRows:
		return transfer.WalletDTO{}, nil
	default:
		return walletInDB.ToDTO(), err
	}
}
