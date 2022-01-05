package testdb

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/skwol/wallet/pkg/client/pgdb"
)

func DBClient() (*pgdb.PGDB, error) {
	client, err := pgdb.NewClient("test")
	if err != nil {
		return nil, fmt.Errorf("error creating db client: %w", err)
	}
	if client == nil {
		return nil, fmt.Errorf("missing db client")
	}
	driver, err := postgres.WithInstance(client.Conn, &postgres.Config{})
	if err != nil {
		return nil, fmt.Errorf("error getting migrate driver: %w", err)
	}
	migrationsPath := "/go/src/github.com/skwol/wallet/pkg/migrations"
	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file:%s", migrationsPath), "postgres", driver)
	if err != nil {
		return nil, fmt.Errorf("error creating migrate instance: %w", err)
	}
	m.Down()
	if err := m.Up(); err != nil {
		return nil, fmt.Errorf("error running up migrations: %w", err)
	}
	return client, nil
}
