package composites

import (
	"fmt"

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

func NewWalletComposite(db *PgDBComposite) (*WalletComposite, error) {
	if db == nil {
		return nil, fmt.Errorf("missing db composite")
	}
	storage, err := dbwallet.NewStorage(db.client)
	if err != nil {
		return nil, fmt.Errorf("error creating account storage %w", err)
	}
	service, err := domainwallet.NewService(storage)
	if err != nil {
		return nil, fmt.Errorf("error creating account service %w", err)
	}
	handler, err := handlerwallet.NewHandler(service)
	if err != nil {
		return nil, fmt.Errorf("error creating account handler %w", err)
	}
	return &WalletComposite{
		Storage: storage,
		Service: service,
		Handler: handler,
	}, nil
}
