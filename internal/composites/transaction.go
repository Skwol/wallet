package composites

import (
	"github.com/pkg/errors"

	"github.com/skwol/wallet/pkg/logging"

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

func NewTransactionComposite(db *PgDBComposite, logger logging.Logger) (*TransactionComposite, error) {
	if db == nil {
		return nil, errors.New("missing db composite")
	}
	storage, err := dbtransaction.NewStorage(db.client, logger)
	if err != nil {
		return nil, errors.Wrap(err, "error creating transaction storage")
	}
	service, err := domaintransaction.NewService(storage, logger)
	if err != nil {
		return nil, errors.Wrap(err, "error creating transaction service")
	}
	handler, err := handlertransaction.NewHandler(service, logger)
	if err != nil {
		return nil, errors.Wrap(err, "error creating transaction handler")
	}
	return &TransactionComposite{
		Storage: storage,
		Service: service,
		Handler: handler,
	}, nil
}
