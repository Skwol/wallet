package transaction

import (
	"context"
	"database/sql"
	"time"

	"github.com/skwol/wallet/internal/domain/transaction"
	"github.com/skwol/wallet/pkg/client/pgdb"
)

type dbTransaction struct {
	ID         int64
	SenderID   int64
	ReceiverID int64
	Amount     float64
	Timestamp  time.Time
	Type       transaction.TranType
}

func (db dbTransaction) ToDTO() transaction.TransactionDTO {
	return transaction.TransactionDTO{
		ID:         db.ID,
		SenderID:   db.SenderID,
		ReceiverID: db.ReceiverID,
		Amount:     db.Amount,
		Timestamp:  db.Timestamp,
		Type:       db.Type,
	}
}

type transactionStorage struct {
	db *pgdb.PGDB
}

func NewStorage(db *pgdb.PGDB) (transaction.Storage, error) {
	return &transactionStorage{db: db}, nil
}

func (as *transactionStorage) GetByID(ctx context.Context, id int64) (transaction.TransactionDTO, error) {
	row := as.db.Conn.QueryRowContext(ctx, "SELECT id, sender_id, receiver_id, amount, date, tran_type FROM transaction WHERE id = $1;", id)
	var tran dbTransaction
	switch err := row.Scan(&tran.ID, &tran.SenderID, &tran.ReceiverID, &tran.Amount, &tran.Timestamp, &tran.Type); err {
	case sql.ErrNoRows:
		return transaction.TransactionDTO{}, nil
	default:
		return tran.ToDTO(), err
	}
}

func (as *transactionStorage) GetAll(ctx context.Context, limit int, offset int) ([]transaction.TransactionDTO, error) {
	var list []transaction.TransactionDTO

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

func (as *transactionStorage) GetFiltered(ctx context.Context, filter *transaction.FilterTransactionsDTO, limit, offset int) ([]transaction.TransactionDTO, error) {
	var list []transaction.TransactionDTO
	dbFilter := newTransactionFilter(filter)
	rows, err := as.db.Conn.QueryContext(ctx, dbFilter.BuildQuery(limit, offset))
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
