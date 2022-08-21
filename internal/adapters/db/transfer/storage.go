package transfer

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"github.com/skwol/wallet/pkg/client/pgdb"
	"github.com/skwol/wallet/pkg/logging"

	"github.com/skwol/wallet/internal/domain/transfer"
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
	db     *pgdb.PGDB
	logger logging.Logger
}

func NewStorage(db *pgdb.PGDB, logger logging.Logger) (transfer.Storage, error) {
	return &transferStorage{db: db, logger: logger}, nil
}

func (ts transferStorage) Create(ctx context.Context, dto *transfer.CreateTransferDTO) (transfer.DTO, error) {
	var result transfer.DTO
	result.CreateTransferDTO = *dto
	tx, err := ts.db.Conn.BeginTx(ctx, nil)
	rollback := func() {
		err := tx.Rollback()
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			ts.logger.Errorf("rollback transaction %s", err)
		}
	}
	if err != nil {
		rollback()
		return result, errors.Wrap(err, "error beginning transaction")
	}
	if _, err = tx.ExecContext(ctx, "UPDATE wallet SET balance=$1 WHERE id=$2;", dto.Sender.Balance, dto.Sender.ID); err != nil {
		rollback()
		return result, errors.Wrap(err, "error updating sender wallet")
	}
	if _, err = tx.ExecContext(ctx, "UPDATE wallet SET balance=$1 WHERE id=$2;", dto.Receiver.Balance, dto.Receiver.ID); err != nil {
		rollback()
		return result, errors.Wrap(err, "error updating receiver wallet")
	}
	row := tx.QueryRow("INSERT INTO transaction (sender_id, receiver_id, amount, date, tran_type) VALUES ($1, $2, $3, $4, 'transfer') RETURNING id;", dto.Sender.ID, dto.Receiver.ID, dto.Amount, dto.Timestamp)

	if err = row.Scan(&result.ID); err != nil {
		rollback()
		return result, err
	}
	if err = tx.Commit(); err != nil {
		return result, errors.Wrap(err, "error during commit")
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
