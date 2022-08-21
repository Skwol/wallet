package composites

import (
	"github.com/pkg/errors"

	"github.com/skwol/wallet/pkg/clock"
	"github.com/skwol/wallet/pkg/logging"

	adapters "github.com/skwol/wallet/internal/adapters/api"
	handlerwallet "github.com/skwol/wallet/internal/adapters/api/wallet"
	dbwallet "github.com/skwol/wallet/internal/adapters/db/wallet"
	domainwallet "github.com/skwol/wallet/internal/domain/wallet"
)

type WalletComposite struct {
	Storage domainwallet.Storage
	Service domainwallet.Service
	Handler adapters.Handler
}

func NewWalletComposite(db *PgDBComposite, logger logging.Logger) (*WalletComposite, error) {
	if db == nil {
		return nil, errors.New("missing db composite")
	}
	storage, err := dbwallet.NewStorage(db.client, logger)
	if err != nil {
		return nil, errors.Wrap(err, "error creating wallet storage")
	}
	service, err := domainwallet.NewService(storage, logger, clock.Real{})
	if err != nil {
		return nil, errors.Wrap(err, "error creating wallet service")
	}
	handler, err := handlerwallet.NewHandler(service, logger)
	if err != nil {
		return nil, errors.Wrap(err, "error creating wallet handler")
	}
	return &WalletComposite{
		Storage: storage,
		Service: service,
		Handler: handler,
	}, nil
}
