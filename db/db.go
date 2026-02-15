package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func GetDb() *pgxpool.Pool {
	if db == nil {
		DbConnect()
	}
	return db
}

func DbConnect() error {
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSW")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", username, password, host, port, dbName)

	config, err := pgxpool.ParseConfig(dbUrl)
	if err != nil {
		return fmt.Errorf("unable to parse database URL: %w", err)
	}

	// config.MaxConns = 10
	// config.MinConns = 2

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	db = pool
	return nil
}
