package composites

import (
	"fmt"

	"github.com/skwol/wallet/pkg/logging"

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

func NewCommonComposite(db *PgDBComposite, logger logging.Logger) (*CommonComposite, error) {
	if db == nil {
		return nil, fmt.Errorf("missing db composite")
	}
	storage, err := dbcommon.NewStorage(db.client, logger)
	if err != nil {
		return nil, fmt.Errorf("error creating common storage %w", err)
	}
	service, err := domaincommon.NewService(storage, logger)
	if err != nil {
		return nil, fmt.Errorf("error creating common service %w", err)
	}
	handler, err := handlercommon.NewHandler(service, logger)
	if err != nil {
		return nil, fmt.Errorf("error creating common handler %w", err)
	}
	return &CommonComposite{
		Storage: storage,
		Service: service,
		Handler: handler,
	}, nil
}
