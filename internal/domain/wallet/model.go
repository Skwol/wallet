package wallet

import (
	"fmt"
	"math"
	"time"

	"github.com/skwol/wallet/internal/domain/transaction"
)

var timeNow = func() time.Time {
	return time.Now()
}

type Wallet struct {
	ID                  int64
	Name                string
	AccountID           int64
	Balance             float64
	TransactionsToApply []transaction.Transaction
}

func (w Wallet) toDTO() *WalletDTO {
	transactionsToApply := make([]*transaction.TransactionDTO, len(w.TransactionsToApply))
	for i, tran := range w.TransactionsToApply {
		transactionsToApply[i] = tran.ToDTO()
	}
	return &WalletDTO{
		ID:                  w.ID,
		Name:                w.Name,
		AccountID:           w.AccountID,
		Balance:             w.Balance,
		TransactionsToApply: transactionsToApply,
	}
}

func (w *Wallet) Update(walletDTO *UpdateWalletDTO) (*Wallet, error) {
	if w.AccountID != walletDTO.AccountID {
		return nil, fmt.Errorf("wrong account")
	}
	if walletDTO.Balance < 0 {
		return nil, fmt.Errorf("balance can not be less then 0")
	}
	if walletDTO.Balance == w.Balance {
		return nil, fmt.Errorf("balance should be updated")
	}
	var tType transaction.TranType
	if walletDTO.Balance > w.Balance {
		tType = transaction.TranTypeDeposit
	} else {
		tType = transaction.TranTypeWithdraw
	}
	w.TransactionsToApply = append(w.TransactionsToApply, transaction.Transaction{
		SenderID:   w.ID,
		ReceiverID: w.ID,
		Amount:     math.Abs(w.Balance - walletDTO.Balance),
		Timestamp:  timeNow(),
		Type:       tType,
	})
	w.Balance = walletDTO.Balance

	return w, nil
}
