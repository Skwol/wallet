package composites

import (
	"fmt"

	adapters "github.com/skwol/wallet/internal/adapters/api"
	handlertransaction "github.com/skwol/wallet/internal/adapters/api/transaction"
	dbtransaction "github.com/skwol/wallet/internal/adapters/db/transaction"
	domaintransaction "github.com/skwol/wallet/internal/domain/transaction"
)

type TransactionComposite struct {
	Storage domaintransaction.Storage
	Service domaintransaction.Service
	Handler adapters.Handler
}

func NewTransactionComposite(db *PgDBComposite) (*TransactionComposite, error) {
	if db == nil {
		return nil, fmt.Errorf("missing db composite")
	}
	storage, err := dbtransaction.NewStorage(db.client)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction storage %w", err)
	}
	service, err := domaintransaction.NewService(storage)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction service %w", err)
	}
	handler, err := handlertransaction.NewHandler(service)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction handler %w", err)
	}
	return &TransactionComposite{
		Storage: storage,
		Service: service,
		Handler: handler,
	}, nil
}
