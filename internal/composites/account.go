package composites

import (
	"fmt"

	adapters "github.com/skwol/wallet/internal/adapters/api"
	handleraccount "github.com/skwol/wallet/internal/adapters/api/account"
	dbacct "github.com/skwol/wallet/internal/adapters/db/account"
	domainacct "github.com/skwol/wallet/internal/domain/account"
)

type AccountComposite struct {
	Storage domainacct.Storage
	Service domainacct.Service
	Handler adapters.Handler
}

func NewAccountComposite(db *PgDBComposite) (*AccountComposite, error) {
	if db == nil {
		return nil, fmt.Errorf("missing db composite")
	}
	storage, err := dbacct.NewStorage(db.client)
	if err != nil {
		return nil, fmt.Errorf("error creating account storage %w", err)
	}
	service, err := domainacct.NewService(storage)
	if err != nil {
		return nil, fmt.Errorf("error creating account service %w", err)
	}
	handler, err := handleraccount.NewHandler(service)
	if err != nil {
		return nil, fmt.Errorf("error creating account handler %w", err)
	}
	return &AccountComposite{
		Storage: storage,
		Service: service,
		Handler: handler,
	}, nil
}
