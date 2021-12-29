package composites

import (
	"fmt"

	adapters "github.com/skwol/wallet/internal/adapters/api"
	handlercommon "github.com/skwol/wallet/internal/adapters/api/common"
	dbwallet "github.com/skwol/wallet/internal/adapters/db/common"
	domaincommon "github.com/skwol/wallet/internal/domain/common"
)

type CommonComposite struct {
	Storage domaincommon.Storage
	Service domaincommon.Service
	Handler adapters.Handler
}

func NewCommonComposite(db *PgDBComposite) (*CommonComposite, error) {
	if db == nil {
		return nil, fmt.Errorf("missing db composite")
	}
	storage, err := dbwallet.NewStorage(db.client)
	if err != nil {
		return nil, fmt.Errorf("error creating account storage %w", err)
	}
	service, err := domaincommon.NewService(storage)
	if err != nil {
		return nil, fmt.Errorf("error creating account service %w", err)
	}
	handler, err := handlercommon.NewHandler(service)
	if err != nil {
		return nil, fmt.Errorf("error creating account handler %w", err)
	}
	return &CommonComposite{
		Storage: storage,
		Service: service,
		Handler: handler,
	}, nil
}
