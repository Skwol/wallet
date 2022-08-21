// Package wallet provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.11.0 DO NOT EDIT.
package wallet

import (
	"time"
)

// Defines values for TransactionType.
const (
	Deposit  TransactionType = "deposit"
	Transfer TransactionType = "transfer"
	Withdraw TransactionType = "withdraw"
)

// Error defines model for Error.
type Error struct {
	Code      *int    `json:"code,omitempty"`
	Error     string  `json:"error"`
	ErrorType *string `json:"errorType,omitempty"`
	Status    string  `json:"status"`
}

// Transaction defines model for Transaction.
type Transaction struct {
	// transfer amount
	Amount *float32 `json:"amount,omitempty"`

	// transaction id
	Id *int `json:"id,omitempty"`

	// receiver wallet id
	ReceiverId *int `json:"receiver_id,omitempty"`

	// sender wallet id
	SenderId  *int             `json:"sender_id,omitempty"`
	Timestamp *time.Time       `json:"timestamp,omitempty"`
	Type      *TransactionType `json:"type,omitempty"`
}

// TransactionType defines model for Transaction.Type.
type TransactionType string

// Wallet defines model for Wallet.
type Wallet struct {
	// Wallet balance
	Balance float32 `json:"balance"`

	// Wallet id
	Id int `json:"id"`

	// Wallet name
	Name         string         `json:"name"`
	Transactions *[]Transaction `json:"transactions,omitempty"`
}

// Wallets defines model for Wallets.
type Wallets struct {
	Wallets *[]Wallet `json:"Wallets,omitempty"`
}

// PathParamWalletID defines model for PathParamWalletID.
type PathParamWalletID = float32

// QueryParamLimit defines model for QueryParamLimit.
type QueryParamLimit = float32

// QueryParamOffset defines model for QueryParamOffset.
type QueryParamOffset = float32

// GetWalletsParams defines parameters for GetWallets.
type GetWalletsParams struct {
	// Limit of how many records returned
	Limit QueryParamLimit `form:"limit" json:"limit"`

	// Offset of returned records
	Offset QueryParamOffset `form:"offset" json:"offset"`
}

// GetWalletWithTransactionsParams defines parameters for GetWalletWithTransactions.
type GetWalletWithTransactionsParams struct {
	// Limit of how many records returned
	Limit QueryParamLimit `form:"limit" json:"limit"`

	// Offset of returned records
	Offset QueryParamOffset `form:"offset" json:"offset"`
}