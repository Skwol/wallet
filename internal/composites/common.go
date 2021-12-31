package composites

import (
	"fmt"

	adapters "github.com/skwol/wallet/internal/adapters/api"
	handlercommon "github.com/skwol/wallet/internal/adapters/api/common"
	dbcommon "github.com/skwol/wallet/internal/adapters/db/common"
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
	storage, err := dbcommon.NewStorage(db.client)
	if err != nil {
		return nil, fmt.Errorf("error creating common storage %w", err)
	}
	service, err := domaincommon.NewService(storage)
	if err != nil {
		return nil, fmt.Errorf("error creating common service %w", err)
	}
	handler, err := handlercommon.NewHandler(service)
	if err != nil {
		return nil, fmt.Errorf("error creating common handler %w", err)
	}
	return &CommonComposite{
		Storage: storage,
		Service: service,
		Handler: handler,
	}, nil
}
