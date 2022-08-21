package composites

import (
	"context"
	"fmt"

	"github.com/skwol/wallet/pkg/client/pgdb"
)

type PgDBComposite struct {
	client *pgdb.PGDB
}

func NewPgDBComposite(ctx context.Context) (*PgDBComposite, error) {
	client, err := pgdb.NewClient("production")
	if err != nil {
		return nil, fmt.Errorf("error getting psql client: %w", err)
	}
	if client == nil {
		return nil, fmt.Errorf("missing psql client")
	}
	return &PgDBComposite{client: client}, nil
}
