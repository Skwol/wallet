package composites

import (
	"context"

	"github.com/pkg/errors"

	"github.com/skwol/wallet/pkg/client/pgdb"
)

type PgDBComposite struct {
	client *pgdb.PGDB
}

func NewPgDBComposite(ctx context.Context) (*PgDBComposite, error) {
	client, err := pgdb.NewClient("production")
	if err != nil {
		return nil, errors.Wrap(err, "error getting psql client")
	}
	if client == nil {
		return nil, errors.New("missing psql client")
	}
	return &PgDBComposite{client: client}, nil
}
