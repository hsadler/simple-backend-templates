package database

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pashagolub/pgxmock/v3"
)

type PgxPoolIface interface {
	Begin(context.Context) (pgx.Tx, error)
	Exec(context.Context, string, ...interface{}) (pgconn.CommandTag, error)
	QueryRow(context.Context, string, ...interface{}) pgx.Row
	Query(context.Context, string, ...interface{}) (pgx.Rows, error)
	Ping(context.Context) error
	Close()
}

func SetupDB() (PgxPoolIface, error) {
	// Connect to database and create tables
	var dbpool *pgxpool.Pool
	var err error
	dbpool, err = pgxpool.New(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	log.Println("Connected to database")
	CreateTables(dbpool)
	log.Println("Created tables")
	return dbpool, err
}

func SetupTestDB() (PgxPoolIface, error) {
	return pgxmock.NewPool()
}

func CreateTables(dbpool *pgxpool.Pool) {
	_, err := dbpool.Exec(context.Background(), `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		CREATE TABLE IF NOT EXISTS item (
			id SERIAL PRIMARY KEY,
			uuid UUID DEFAULT uuid_generate_v4(),
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			name VARCHAR(50),
			price NUMERIC(10, 2),
			CONSTRAINT name_unique UNIQUE (name)
		);
	`)
	if err != nil {
		log.Fatal(err)
	}
}
