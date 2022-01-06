package common

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
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

func (cs *commonStorage) GenerateFakeData(ctx context.Context, numberOfRecordsToCreate int) error {
	rand.Seed(time.Now().UnixNano())
	var wg sync.WaitGroup

	generateData := func(ctx context.Context, start, end int) error {
		defer wg.Done()
		for i := start; i <= end; i++ {
			select {
			case <-ctx.Done():
				return nil
			default:
				var walletID int
				walletName := fmt.Sprintf("wallet_%d", i+1)
				walletBalance := randFloat(1, 1200)

				row := cs.db.Conn.QueryRow("INSERT INTO wallet (name, balance) VALUES ($1, $2) RETURNING id;", walletName, walletBalance)
				if err := row.Scan(&walletID); err != nil {
					return err
				}

				if _, err := cs.db.Conn.ExecContext(ctx, "INSERT INTO transaction (sender_id, receiver_id, amount, date, tran_type) VALUES ($1, $1, $2, current_timestamp, $3);", walletID, walletBalance, transaction.TranTypeDeposit); err != nil {
					return err
				}
			}
		}
		return nil
	}

	ctxForData, cancel := context.WithTimeout(context.Background(), time.Minute*25)
	defer cancel()

	numberOfBatches := 20
	recPerBatch := numberOfRecordsToCreate / numberOfBatches

	for i := 1; i <= numberOfBatches; i++ {
		wg.Add(1)
		go generateData(ctxForData, recPerBatch*i-recPerBatch+1, recPerBatch*i)
	}
	wg.Wait()
	return nil
}

func randFloat(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
