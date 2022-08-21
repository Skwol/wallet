package wallet

import (
	"context"
	"database/sql"
	"time"

	"github.com/pkg/errors"

	"github.com/skwol/wallet/pkg/client/pgdb"
	"github.com/skwol/wallet/pkg/logging"

	"github.com/skwol/wallet/internal/domain/wallet"
)

type dbWallet struct {
	ID      int64
	Name    string
	Balance float64
}

func (db dbWallet) ToDTO() wallet.DTO {
	return wallet.DTO{
		ID:      db.ID,
		Name:    db.Name,
		Balance: db.Balance,
	}
}

type dbTransaction struct {
	ID         int64
	SenderID   int64
	ReceiverID int64
	Amount     float64
	Timestamp  time.Time
	Type       wallet.TranType
}

func (db dbTransaction) ToDTO() wallet.TransactionDTO {
	return wallet.TransactionDTO{
		ID:         db.ID,
		SenderID:   db.SenderID,
		ReceiverID: db.ReceiverID,
		Amount:     db.Amount,
		Timestamp:  db.Timestamp,
		Type:       db.Type,
	}
}

type walletStorage struct {
	db     *pgdb.PGDB
	logger logging.Logger
}

func NewStorage(db *pgdb.PGDB, logger logging.Logger) (wallet.Storage, error) {
	return &walletStorage{db: db, logger: logger}, nil
}

func (as *walletStorage) Create(ctx context.Context, dto wallet.DTO) (wallet.DTO, error) {
	tx, err := as.db.Conn.BeginTx(ctx, nil)
	if err != nil {
		return dto, err
	}

	row := tx.QueryRowContext(ctx, "INSERT INTO wallet (name, balance) VALUES ($1, $2) RETURNING id;", dto.Name, dto.Balance)
	rollback := func() {
		err := tx.Rollback()
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			as.logger.Errorf("rollback transaction %s", err)
		}
	}
	if err = row.Scan(&dto.ID); err != nil {
		rollback()
		return dto, err
	}

	for _, tran := range dto.TransactionsToApply {
		_, err = tx.ExecContext(ctx, "INSERT INTO transaction (sender_id, receiver_id, amount, date, tran_type) VALUES ($1, $1, $2, $3, $4);", dto.ID, tran.Amount, tran.Timestamp, tran.Type)
		if err != nil {
			rollback()
			return dto, errors.Wrap(err, "error inserting transaction")
		}
	}

	if err := tx.Commit(); err != nil {
		return dto, errors.Wrap(err, "error committing transction")
	}

	return dto, nil
}

func (as *walletStorage) GetByID(ctx context.Context, id int64) (wallet.DTO, error) {
	query := `SELECT id, name, balance FROM wallet WHERE id = $1;`
	row := as.db.Conn.QueryRowContext(ctx, query, id)
	var walletInDB dbWallet
	switch err := row.Scan(&walletInDB.ID, &walletInDB.Name, &walletInDB.Balance); err {
	case sql.ErrNoRows:
		return wallet.DTO{}, nil
	default:
		return walletInDB.ToDTO(), err
	}
}

func (as *walletStorage) GetByIDWithTransactions(ctx context.Context, id int64, limit int, offset int) (wallet.DTO, error) {
	tx, err := as.db.Conn.BeginTx(ctx, nil)
	if err != nil {
		return wallet.DTO{}, errors.Wrap(err, "error beginning transaction")
	}
	query := `SELECT id, name, balance FROM wallet WHERE id = $1;`
	row := as.db.Conn.QueryRowContext(ctx, query, id)
	var walletInDB dbWallet
	if err := row.Scan(&walletInDB.ID, &walletInDB.Name, &walletInDB.Balance); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return wallet.DTO{}, nil
		}
		return wallet.DTO{}, err
	}

	query = "SELECT id, sender_id, receiver_id, amount, date, tran_type FROM transaction WHERE sender_id = $1 OR receiver_id = $1 ORDER BY ID ASC LIMIT $2 OFFSET $3"
	rows, err := as.db.Conn.Query(query, walletInDB.ID, limit, offset)
	if err != nil {
		return wallet.DTO{}, err
	}
	var (
		list []wallet.TransactionDTO
		tran dbTransaction
	)
	for rows.Next() {
		if err := rows.Scan(&tran.ID, &tran.SenderID, &tran.ReceiverID, &tran.Amount, &tran.Timestamp, &tran.Type); err != nil {
			return wallet.DTO{}, err
		}
		list = append(list, tran.ToDTO())
	}

	if err := tx.Commit(); err != nil {
		return wallet.DTO{}, errors.Wrap(err, "error committing transction")
	}
	walletDTO := walletInDB.ToDTO()
	walletDTO.Transactions = list
	return walletDTO, nil
}

func (as *walletStorage) GetByName(ctx context.Context, name string) (wallet.DTO, error) {
	query := `SELECT id FROM wallet WHERE name = $1;`
	row := as.db.Conn.QueryRowContext(ctx, query, name)
	var walletInDB dbWallet
	switch err := row.Scan(&walletInDB.ID); err {
	case sql.ErrNoRows:
		return wallet.DTO{}, nil
	default:
		return walletInDB.ToDTO(), err
	}
}

func (as *walletStorage) GetAll(ctx context.Context, limit int, offset int) ([]wallet.DTO, error) {
	var list []wallet.DTO
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

func (as *walletStorage) Update(ctx context.Context, walletDTO wallet.DTO) error {
	tx, err := as.db.Conn.BeginTx(ctx, nil)
	if err != nil {
		return errors.Wrap(err, "error beginning transaction")
	}

	rollback := func() {
		err := tx.Rollback()
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			as.logger.Errorf("rollback transaction %s", err)
		}
	}
	_, err = tx.ExecContext(ctx, "UPDATE wallet SET name=$1, balance=$2 WHERE id=$3;", walletDTO.Name, walletDTO.Balance, walletDTO.ID)
	if err != nil {
		rollback()
		return errors.Wrap(err, "error updating wallet")
	}

	for _, tran := range walletDTO.TransactionsToApply {
		_, err = tx.ExecContext(ctx, "INSERT INTO transaction (sender_id, receiver_id, amount, date, tran_type) VALUES ($1, $1, $2, $3, $4);", walletDTO.ID, tran.Amount, tran.Timestamp, tran.Type)
		if err != nil {
			rollback()
			return errors.Wrap(err, "error inserting transaction")
		}
	}

	if err := tx.Commit(); err != nil {
		return errors.Wrap(err, "error committing transction")
	}

	return nil
}
