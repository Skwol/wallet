package composites

import (
	"github.com/pkg/errors"

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
		return nil, errors.New("missing db composite")
	}
	storage, err := dbcommon.NewStorage(db.client, logger)
	if err != nil {
		return nil, errors.Wrap(err, "error creating common storage")
	}
	service, err := domaincommon.NewService(storage, logger)
	if err != nil {
		return nil, errors.Wrap(err, "error creating common service")
	}
	handler, err := handlercommon.NewHandler(service, logger)
	if err != nil {
		return nil, errors.Wrap(err, "error creating common handler")
	}
	return &CommonComposite{
		Storage: storage,
		Service: service,
		Handler: handler,
	}, nil
}
