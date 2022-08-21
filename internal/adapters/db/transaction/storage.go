package transaction

import (
	"context"
	"database/sql"
	"time"

	"github.com/skwol/wallet/pkg/client/pgdb"
	"github.com/skwol/wallet/pkg/logging"

	"github.com/skwol/wallet/internal/domain/transaction"
)

type dbTransaction struct {
	ID         int64
	SenderID   int64
	ReceiverID int64
	Amount     float64
	Timestamp  time.Time
	Type       transaction.TranType
}

func (db dbTransaction) ToDTO() transaction.DTO {
	return transaction.DTO{
		ID:         db.ID,
		SenderID:   db.SenderID,
		ReceiverID: db.ReceiverID,
		Amount:     db.Amount,
		Timestamp:  db.Timestamp,
		Type:       db.Type,
	}
}

type transactionStorage struct {
	db     *pgdb.PGDB
	logger logging.Logger
}

func NewStorage(db *pgdb.PGDB, logger logging.Logger) (transaction.Storage, error) {
	return &transactionStorage{db: db, logger: logger}, nil
}

func (as *transactionStorage) GetByID(ctx context.Context, id int64) (transaction.DTO, error) {
	row := as.db.Conn.QueryRowContext(ctx, "SELECT id, sender_id, receiver_id, amount, date, tran_type FROM transaction WHERE id = $1;", id)
	var tran dbTransaction
	switch err := row.Scan(&tran.ID, &tran.SenderID, &tran.ReceiverID, &tran.Amount, &tran.Timestamp, &tran.Type); err {
	case sql.ErrNoRows:
		return transaction.DTO{}, nil
	default:
		return tran.ToDTO(), err
	}
}

func (as *transactionStorage) GetAll(ctx context.Context, limit int, offset int) ([]transaction.DTO, error) {
	var list []transaction.DTO

	rows, err := as.db.Conn.QueryContext(ctx, "SELECT id, sender_id, receiver_id, amount, date, tran_type FROM transaction ORDER BY ID ASC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return list, err
	}
	var tran dbTransaction
	for rows.Next() {
		if err := rows.Scan(&tran.ID, &tran.SenderID, &tran.ReceiverID, &tran.Amount, &tran.Timestamp, &tran.Type); err != nil {
			return nil, err
		}
		list = append(list, tran.ToDTO())
	}
	return list, nil
}

func (as *transactionStorage) GetFiltered(ctx context.Context, filter *transaction.FilterTransactionsDTO, limit, offset int) ([]transaction.DTO, error) {
	var list []transaction.DTO
	rows, err := as.db.Conn.QueryContext(ctx, newTransactionFilter(filter).BuildQuery(limit, offset))
	if err != nil {
		return list, err
	}
	for rows.Next() {
		var transaction dbTransaction
		err := rows.Scan(&transaction.ID, &transaction.SenderID, &transaction.ReceiverID, &transaction.Amount, &transaction.Timestamp, &transaction.Type)
		if err != nil {
			return list, err
		}
		list = append(list, transaction.ToDTO())
	}
	return list, nil
}
