package wallet

import (
	"math"
	"time"

	"github.com/pkg/errors"
)

const (
	TranTypeDeposit  TranType = "deposit"
	TranTypeWithdraw TranType = "withdraw"
	TranTypeTransfer TranType = "transfer"
)

var (
	ErrNegativeBalance            = errors.New("balance can not be less then 0")
	ErrMissingName                = errors.New("wallet must have a name")
	ErrUpdateWithoutBalanceChange = errors.New("balance must be updated")
)

type TranType string

// Wallet is an aggregate entity
type Wallet struct {
	ID                  int64
	Name                string
	Balance             float64
	TransactionsToApply []Transaction
}

func newWallet(dto *CreateWalletDTO, timestamp time.Time) (*Wallet, error) {
	if err := dto.validate(); err != nil {
		return nil, err
	}
	var transactionsToApply []Transaction
	if dto.Balance > 0 {
		transactionsToApply = append(transactionsToApply, Transaction{Amount: dto.Balance, Timestamp: timestamp, Type: TranTypeDeposit})
	}
	return &Wallet{
		Name:                dto.Name,
		Balance:             dto.Balance,
		TransactionsToApply: transactionsToApply,
	}, nil
}

func (w Wallet) toDTO() DTO {
	transactionsToApply := make([]TransactionDTO, len(w.TransactionsToApply))
	for i, tran := range w.TransactionsToApply {
		transactionsToApply[i] = tran.toDTO()
	}
	return DTO{
		ID:                  w.ID,
		Name:                w.Name,
		Balance:             w.Balance,
		TransactionsToApply: transactionsToApply,
	}
}

func (w *Wallet) Update(walletDTO *UpdateWalletDTO, timestamp time.Time) (*Wallet, error) {
	if err := walletDTO.validate(); err != nil {
		return nil, err
	}
	if walletDTO.Balance == w.Balance {
		return nil, ErrUpdateWithoutBalanceChange
	}
	var tType TranType
	if walletDTO.Balance > w.Balance {
		tType = TranTypeDeposit
	} else {
		tType = TranTypeWithdraw
	}
	w.TransactionsToApply = append(w.TransactionsToApply, Transaction{
		SenderID:   w.ID,
		ReceiverID: w.ID,
		Amount:     math.Abs(w.Balance - walletDTO.Balance),
		Timestamp:  timestamp,
		Type:       tType,
	})
	w.Balance = walletDTO.Balance

	return w, nil
}

type Transaction struct {
	ID         int64
	SenderID   int64
	ReceiverID int64
	Amount     float64
	Timestamp  time.Time
	Type       TranType
}

func (t Transaction) toDTO() TransactionDTO {
	return TransactionDTO(t)
}
