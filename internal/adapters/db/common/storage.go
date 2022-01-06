package common

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/skwol/wallet/internal/domain/common"
	"github.com/skwol/wallet/internal/domain/transaction"
	"github.com/skwol/wallet/pkg/client/pgdb"
)

type commonStorage struct {
	db *pgdb.PGDB
}

func NewStorage(db *pgdb.PGDB) (common.Storage, error) {
	return &commonStorage{db: db}, nil
}

func (cs *commonStorage) GenerateFakeData(ctx context.Context) error {
	tx, err := cs.db.Conn.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("error beginning transaction")
	}

	rand.Seed(time.Now().UnixNano())
	balances := randFloats(0, 1200, 100)

	var walletID int64
	for i, walletBalance := range balances {
		walletName := fmt.Sprintf("wallet_%d", i+1)

		row := tx.QueryRow("INSERT INTO wallet (name, balance) VALUES ($1, $2) RETURNING id;", walletName, walletBalance)
		if err = row.Scan(&walletID); err != nil {
			tx.Rollback()
			return err
		}
		_, err = tx.ExecContext(ctx, "INSERT INTO transaction (sender_id, receiver_id, amount, date, tran_type) VALUES ($1, $1, $2, current_timestamp, $3);", walletID, walletBalance, transaction.TranTypeDeposit)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("error commiting transaction, rolled back")
	}
	return nil
}

func randFloats(min, max float64, n int) []float64 {
	res := make([]float64, n)
	for i := range res {
		res[i] = min + rand.Float64()*(max-min)
	}
	return res
}
