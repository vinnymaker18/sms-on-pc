package storage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	dbURL = "postgres://smsapp@localhost:5432/smsdb"
)

var (
	pool            *pgxpool.Pool
	dbPoolInitError error
)

func init() {
	pool, dbPoolInitError = pgxpool.Connect(context.Background(), dbURL)
	if dbPoolInitError != nil {
		log.Fatal("Error setting up a connection pool")
	}
}

func newDBConn() (*pgxpool.Conn, error) {
	return pool.Acquire(context.Background())
}
