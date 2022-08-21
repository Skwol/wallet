package composites

import (
	"github.com/pkg/errors"

	"github.com/skwol/wallet/pkg/clock"
	"github.com/skwol/wallet/pkg/logging"

	adapters "github.com/skwol/wallet/internal/adapters/api"
	handlertransfer "github.com/skwol/wallet/internal/adapters/api/transfer"
	dbtransfer "github.com/skwol/wallet/internal/adapters/db/transfer"
	domaintransfer "github.com/skwol/wallet/internal/domain/transfer"
)

type TransferComposite struct {
	Storage domaintransfer.Storage
	Service domaintransfer.Service
	Handler adapters.Handler
}

func NewTransferComposite(db *PgDBComposite, logger logging.Logger, clk clock.Clock) (*TransferComposite, error) {
	if db == nil {
		return nil, errors.New("missing db composite")
	}
	storage, err := dbtransfer.NewStorage(db.client, logger)
	if err != nil {
		return nil, errors.Wrap(err, "error creating transaction storage")
	}
	service, err := domaintransfer.NewService(storage, logger, clk)
	if err != nil {
		return nil, errors.Wrap(err, "error creating transaction service")
	}
	handler, err := handlertransfer.NewHandler(service, logger)
	if err != nil {
		return nil, errors.Wrap(err, "error creating transaction handler")
	}
	return &TransferComposite{
		Storage: storage,
		Service: service,
		Handler: handler,
	}, nil
}
