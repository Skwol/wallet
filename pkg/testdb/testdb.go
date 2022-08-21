package testdb

import (
	"fmt"

	migrate "github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/pkg/errors"

	"github.com/skwol/wallet/pkg/client/pgdb"
	"github.com/skwol/wallet/pkg/logging"
)

func DBClient(logger logging.Logger) (*pgdb.PGDB, error) {
	client, err := pgdb.NewClient("test")
	if err != nil {
		return nil, errors.Wrap(err, "error creating db client")
	}
	if client == nil {
		return nil, errors.New("missing db client")
	}
	driver, err := postgres.WithInstance(client.Conn, &postgres.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "error getting migrate driver")
	}
	migrationsPath := "/go/src/github.com/skwol/wallet/db/migrations"
	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file:%s", migrationsPath), "postgres", driver)
	if err != nil {
		return nil, errors.Wrap(err, "error creating migrate instance")
	}
	if err := m.Down(); err != nil {
		logger.Warn("error during migrations down %s", err)
	}
	if err := m.Up(); err != nil {
		return nil, errors.Wrap(err, "error running up migrations")
	}
	return client, nil
}
