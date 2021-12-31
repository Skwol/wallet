package transaction

import (
	"context"
	"database/sql"
	"time"

	"github.com/skwol/wallet/internal/domain/transaction"
	"github.com/skwol/wallet/pkg/client/pgdb"
)

type dbTransaction struct {
	ID         int64                `json:"id,omitempty"`
	SenderID   int64                `json:"sender_id,omitempty"`
	ReceiverID int64                `json:"receiver_id,omitempty"`
	Amount     float64              `json:"amount,omitempty"`
	Timestamp  time.Time            `json:"timestamp,omitempty"`
	Type       transaction.TranType `json:"type,omitempty"`
}

func (db dbTransaction) ToDTO() *transaction.TransactionDTO {
	return &transaction.TransactionDTO{
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

func (as *transactionStorage) Create(ctx context.Context, acct *transaction.TransactionDTO) (*transaction.TransactionDTO, error) {
	return nil, nil
}

func (as *transactionStorage) GetByID(ctx context.Context, id int64) (*transaction.TransactionDTO, error) {
	row := as.db.Conn.QueryRow(`SELECT id, sender_id, receiver_id, amount, date, tran_type FROM transaction WHERE id = $1;`, id)
	var tran dbTransaction
	switch err := row.Scan(&tran.ID, &tran.SenderID, &tran.ReceiverID, &tran.Amount, &tran.Timestamp, &tran.Type); err {
	case sql.ErrNoRows:
		return nil, nil
	default:
		return tran.ToDTO(), err
	}
}

func (as *transactionStorage) GetAll(ctx context.Context, limit int, offset int) ([]*transaction.TransactionDTO, error) {
	var list []*transaction.TransactionDTO

	rows, err := as.db.Conn.Query("SELECT id, sender_id, receiver_id, amount, date, tran_type FROM transaction ORDER BY ID ASC LIMIT $1 OFFSET $2", limit, offset)
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
