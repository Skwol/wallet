package common

import (
	"context"
	"fmt"

	"github.com/Pallinder/go-randomdata"
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

	numberOfAccounts := 100
	accountNames := make(map[string]bool, numberOfAccounts)
	walletNames := make(map[string]bool, numberOfAccounts)

	var accountID, walletID int64
	for i := 0; i < numberOfAccounts; i++ {
		accountName := randomdata.FullName(randomdata.RandomGender)
		walletName := randomdata.SillyName()
		if accountNames[accountName] {
			i--
			continue
		}
		if walletNames[walletName] {
			i--
			continue
		}
		accountNames[accountName] = true
		walletNames[walletName] = true
		walletBalance := randomdata.Decimal(1000)

		row := tx.QueryRowContext(ctx, "INSERT INTO account (username) VALUES ($1)  RETURNING id;", accountName)
		if err := row.Scan(&accountID); err != nil {
			tx.Rollback()
			return err
		}

		row = tx.QueryRow("INSERT INTO wallet (name, account_id, balance) VALUES ($1, $2, $3) RETURNING id;", walletName, accountID, walletBalance)
		err = row.Scan(&walletID)
		if err != nil {
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
