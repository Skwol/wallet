package common

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/skwol/wallet/pkg/client/pgdb"
	"github.com/skwol/wallet/pkg/logging"

	"github.com/skwol/wallet/internal/domain/common"
	"github.com/skwol/wallet/internal/domain/transaction"
)

type commonStorage struct {
	db     *pgdb.PGDB
	logger logging.Logger
}

func NewStorage(db *pgdb.PGDB, logger logging.Logger) (common.Storage, error) {
	return &commonStorage{db: db, logger: logger}, nil
}

func (cs *commonStorage) GenerateFakeData(ctx context.Context, numberOfRecordsToCreate int) error {
	var wg sync.WaitGroup

	generateData := func(ctx context.Context, start, end int) {
		defer wg.Done()
		for i := start; i <= end; i++ {
			select {
			case <-ctx.Done():
				return
			default:
				var walletID int
				walletName := fmt.Sprintf("wallet_%d", i+1)
				walletBalance := randFloat(1, 1200)

				row := cs.db.Conn.QueryRow("INSERT INTO wallet (name, balance) VALUES ($1, $2) RETURNING id;", walletName, walletBalance)
				if err := row.Scan(&walletID); err != nil {
					cs.logger.Warn("error receiving walletID %s", err.Error())
					return
				}

				if _, err := cs.db.Conn.ExecContext(ctx, "INSERT INTO transaction (sender_id, receiver_id, amount, date, tran_type) VALUES ($1, $1, $2, current_timestamp, $3);", walletID, walletBalance, transaction.TranTypeDeposit); err != nil {
					cs.logger.Warn("error inserting transaction %s", err.Error())
					return
				}
			}
		}
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

// Float64 is a shortcut for generating a random float between 0 and
// 1 using crypto/rand.
func randFloat(min, max float64) float64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(1<<53))
	if err != nil {
		return 0
	}

	randomFloat := float64(nBig.Int64()) / (1 << 53)
	return min + randomFloat*(max-min)
}
