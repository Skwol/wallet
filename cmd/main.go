package main

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/skwol/wallet/internal/composites"
	"github.com/skwol/wallet/pkg/logging"
)

func main() {
	logging.Init()
	logger := logging.GetLogger()
	ctx := context.Background()

	logger.Info("router initialization")
	router := mux.NewRouter()

	logger.Info("create db composite")
	db, err := composites.NewPgDBComposite(ctx)
	if err != nil {
		logger.Fatal("db composite failed: ", err.Error())
	}

	logger.Info("create account composite")
	accountComposite, err := composites.NewAccountComposite(db)
	if err != nil {
		logger.Fatal("account composite failed:", err.Error())
	}
	accountComposite.Handler.Register(router)

	logger.Info("create transaction composite")
	transactionComposite, err := composites.NewTransactionComposite(db)
	if err != nil {
		logger.Fatal("transaction composite failed:", err.Error())
	}
	transactionComposite.Handler.Register(router)

	logger.Info("create wallet composite")
	walletComposite, err := composites.NewWalletComposite(db)
	if err != nil {
		logger.Fatal("wallet composite failed:", err.Error())
	}
	walletComposite.Handler.Register(router)

	addr := ":8080"
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatal("Error occurred:", err.Error())
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}
	logger.Fatal(server.Serve(listener))
}
