package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fabianpoels/ai-bubble-tracker/models"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var bunDB *bun.DB

// GetDB returns the Bun DB instance
func GetDB() *bun.DB {
	if bunDB == nil {
		if err := Connect(); err != nil {
			panic(err)
		}
	}
	return bunDB
}

// Connect creates a Bun DB instance
func Connect() error {
	dbUrl := buildDatabaseURL()

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dbUrl)))

	// Configure connection pool
	sqldb.SetMaxOpenConns(10)
	sqldb.SetMaxIdleConns(5)
	sqldb.SetConnMaxLifetime(time.Hour)

	bunDB = bun.NewDB(sqldb, pgdialect.New())

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := bunDB.PingContext(ctx); err != nil {
		return fmt.Errorf("unable to ping database: %w", err)
	}

	return nil
}

// buildDatabaseURL constructs the database connection string
func buildDatabaseURL() string {
	username := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	host := os.Getenv("POSTGRES_HOST")
	port := os.Getenv("POSTGRES_PORT")
	dbName := os.Getenv("POSTGRES_DB")

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		username, password, host, port, dbName)
}

// Close closes the database connection
func Close() {
	if bunDB != nil {
		bunDB.Close()
	}
}

func CreateDatabase() error {
	log.Println("Creating database tables...")

	bunDB := GetDB()
	ctx := context.Background()

	modelsToCreate := []interface{}{
		(*models.Datapoint)(nil),
	}

	for _, model := range modelsToCreate {
		_, err := bunDB.NewCreateTable().
			Model(model).
			IfNotExists().
			Exec(ctx)

		if err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	if err := createIndexes(bunDB, ctx); err != nil {
		return fmt.Errorf("failed to create indexes: %w", err)
	}

	log.Println("✓ Database tables created successfully")
	return nil
}

func createIndexes(bunDB *bun.DB, ctx context.Context) error {
	// Add any custom indexes here if needed
	return nil
}

func DropDatabase() error {
	log.Println("Dropping database tables...")

	bunDB := GetDB()
	ctx := context.Background()

	modelsToDrop := []interface{}{
		(*models.Datapoint)(nil),
	}

	for _, model := range modelsToDrop {
		_, err := bunDB.NewDropTable().
			Model(model).
			IfExists().
			Cascade().
			Exec(ctx)

		if err != nil {
			return fmt.Errorf("failed to drop table: %w", err)
		}
	}

	log.Println("✓ Database tables dropped successfully")
	return nil
}
