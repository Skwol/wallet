package composites

import (
	"fmt"

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
		return nil, fmt.Errorf("missing db composite")
	}
	storage, err := dbtransfer.NewStorage(db.client, logger)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction storage %w", err)
	}
	service, err := domaintransfer.NewService(storage, logger, clk)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction service %w", err)
	}
	handler, err := handlertransfer.NewHandler(service)
	if err != nil {
		return nil, fmt.Errorf("error creating transaction handler %w", err)
	}
	return &TransferComposite{
		Storage: storage,
		Service: service,
		Handler: handler,
	}, nil
}
