package wallet

import (
	"github.com/skwol/wallet/internal/domain/wallet"
)

func newWallet(dto wallet.DTO) Wallet {
	w := Wallet{
		Id:      int(dto.ID),
		Name:    dto.Name,
		Balance: float32(dto.Balance),
	}
	if len(dto.Transactions) == 0 {
		return w
	}
	transactions := make([]Transaction, 0, len(dto.Transactions))
	for _, t := range dto.Transactions {
		transactions = append(transactions, newTransaction(t))
	}
	w.Transactions = &transactions
	return w
}

func (w Wallet) toCreateRequest() wallet.CreateWalletDTO {
	return wallet.CreateWalletDTO{
		Name:    w.Name,
		Balance: float64(w.Balance),
	}
}

func (w Wallet) toUpdateRequest() wallet.UpdateWalletDTO {
	return wallet.UpdateWalletDTO{
		CreateWalletDTO: w.toCreateRequest(),
	}
}

func newTransaction(dto wallet.TransactionDTO) Transaction {
	id := int(dto.ID)
	senderID := int(dto.SenderID)
	receiverID := int(dto.ReceiverID)
	amount := float32(dto.Amount)
	tranType := TransactionType(dto.Type)
	return Transaction{
		Id:         &id,
		SenderId:   &senderID,
		ReceiverId: &receiverID,
		Amount:     &amount,
		Timestamp:  &dto.Timestamp,
		Type:       &tranType,
	}
}
