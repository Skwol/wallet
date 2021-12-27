package main

import (
	handleraccount "github.com/skwol/wallet/internal/adapters/api/account"
	handlertransaction "github.com/skwol/wallet/internal/adapters/api/transaction"
	handlerwallet "github.com/skwol/wallet/internal/adapters/api/wallet"
	dbacct "github.com/skwol/wallet/internal/adapters/db/account"
	dbtransaction "github.com/skwol/wallet/internal/adapters/db/transaction"
	dbwallet "github.com/skwol/wallet/internal/adapters/db/wallet"
	domainacct "github.com/skwol/wallet/internal/domain/account"
	domaintransaction "github.com/skwol/wallet/internal/domain/transaction"
	domainwallet "github.com/skwol/wallet/internal/domain/wallet"
)

func main() {
	accountStorage := dbacct.NewStorage()
	accountService := domainacct.NewService(accountStorage)
	accountHandler := handleraccount.NewHandler(accountService)

	transactionStorage := dbtransaction.NewStorage()
	transactionService := domaintransaction.NewService(transactionStorage)
	transactionHandler := handlertransaction.NewHandler(transactionService)

	walletStorage := dbwallet.NewStorage()
	walletService := domainwallet.NewService(walletStorage)
	walletHandler := handlerwallet.NewHandler(walletService)
}
