package transaction

import (
	"context"
	"time"

	"github.com/skwol/wallet/internal/domain/transaction"
	"github.com/skwol/wallet/pkg/client/pgdb"
)

type dbTransaction struct {
	ID              int64                `json:"id,omitempty"`
	SenderID        int64                `json:"sender_id,omitempty"`
	ReceiverID      int64                `json:"receiver_id,omitempty"`
	Amount          float64              `json:"amount,omitempty"`
	Timestamp       time.Time            `json:"timestamp,omitempty"`
	TransactionType transaction.TranType `json:"transaction_type,omitempty"`
}

func (db dbTransaction) ToDTO() *transaction.TransactionDTO {
	return &transaction.TransactionDTO{
		ID:              db.ID,
		SenderID:        db.SenderID,
		ReceiverID:      db.ReceiverID,
		Amount:          db.Amount,
		Timestamp:       db.Timestamp,
		TransactionType: db.TransactionType,
	}
}

type transactionStorage struct {
	db *pgdb.PGDB
}

func NewStorage(db *pgdb.PGDB) (transaction.Storage, error) {
	return &transactionStorage{db: db}, nil
}

func (as *transactionStorage) Create(ctx context.Context, acct *transaction.TransactionDTO) (*transaction.TransactionDTO, error) {
	return nil, nil
}

func (as *transactionStorage) GetByID(ctx context.Context, id int64) (*transaction.TransactionDTO, error) {
	return nil, nil
}

func (as *transactionStorage) GetAll(ctx context.Context, limit int, offset int) ([]*transaction.TransactionDTO, error) {
	var list []*transaction.TransactionDTO
	rows, err := as.db.Conn.Query("SELECT * FROM transaction ORDER BY ID ASC LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		return list, err
	}
	var transaction dbTransaction
	for rows.Next() {
		if err := rows.Scan(&transaction.ID, &transaction.SenderID, &transaction.ReceiverID, &transaction.Amount, &transaction.Timestamp, &transaction.TransactionType); err != nil {
			return nil, err
		}
		list = append(list, transaction.ToDTO())
	}
	return list, nil
}
