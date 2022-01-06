package main

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/skwol/wallet/internal/composites"
	"github.com/skwol/wallet/pkg/logging"
	"github.com/skwol/wallet/pkg/middleware"
)

const addr = ":8080"

func main() {
	logging.Init()
	logger := logging.GetLogger()
	ctx := context.Background()

	logger.Info("router initialization")
	router := mux.NewRouter()
	router.Use(middleware.LimitPerUserByIP)

	logger.Info("create db composite")
	db, err := composites.NewPgDBComposite(ctx)
	if err != nil {
		// sloppy workaround initial container start
		ticker := time.NewTicker(3 * time.Second)
		ctxToReconnect, cancel := context.WithTimeout(ctx, time.Second*12)
	LOOP:
		for {
			select {
			case <-ctxToReconnect.Done():
				cancel()
				break LOOP
			case <-ticker.C:
				db, err = composites.NewPgDBComposite(ctx)
			}
			if err == nil {
				break
			}
		}
		if err != nil {
			logger.Fatal("db composite failed: ", err.Error())
		}
	}

	logger.Info("create transaction composite")
	transactionComposite, err := composites.NewTransactionComposite(db)
	if err != nil {
		logger.Fatal("transaction composite failed:", err.Error())
	}
	transactionComposite.Handler.Register(router)

	logger.Info("create transfer composite")
	transferComposite, err := composites.NewTransferComposite(db)
	if err != nil {
		logger.Fatal("transfer composite failed:", err.Error())
	}
	transferComposite.Handler.Register(router)

	logger.Info("create wallet composite")
	walletComposite, err := composites.NewWalletComposite(db)
	if err != nil {
		logger.Fatal("wallet composite failed:", err.Error())
	}
	walletComposite.Handler.Register(router)

	logger.Info("create common composite")
	commonComposite, err := composites.NewCommonComposite(db)
	if err != nil {
		logger.Fatal("common composite failed:", err.Error())
	}
	commonComposite.Handler.Register(router)

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
