package pgdb

import (
	"database/sql"
	"fmt"
)

const (
	HOST = "database"
	PORT = 5432
)

type PGDB struct {
	Conn *sql.DB
}

func NewClient(username, password, database string) (*PGDB, error) {
	db := &PGDB{}
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		HOST, PORT, username, password, database)
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
