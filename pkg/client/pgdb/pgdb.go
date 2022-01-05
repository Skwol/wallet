package pgdb

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

const (
	HOST      = "database"
	HOST_TEST = "database_test"
	PORT      = 5432
)

type PGDB struct {
	Conn *sql.DB
}

func NewClient(env string) (*PGDB, error) {
	db := &PGDB{}
	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	database := os.Getenv("POSTGRES_DB_PROD")
	host := HOST
	if env != "production" {
		database = os.Getenv("POSTGRES_DB_TEST")
		host = HOST_TEST
	}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, PORT, username, password, database)
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
