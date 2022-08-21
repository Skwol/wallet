package pgdb

import (
	"database/sql"
	"fmt"
	"os"

	// psql connection created in this package
	_ "github.com/lib/pq"
)

const (
	host     = "database"
	hostTest = "database_test"
	port     = 5432
)

type PGDB struct {
	Conn *sql.DB
}

func NewClient(env string) (*PGDB, error) {
	db := &PGDB{}
	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	database := os.Getenv("POSTGRES_DB_PROD")
	host := host
	if env != "production" {
		database = os.Getenv("POSTGRES_DB_TEST")
		host = hostTest
	}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, username, password, database)
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return db, err
	}
	db.Conn = conn
	err = db.Conn.Ping()
	if err != nil {
		return db, err
	}
	return db, nil
}
