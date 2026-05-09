package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/huipizda007/ATB_go/internal/config"
	"github.com/huipizda007/ATB_go/internal/server"
	db "github.com/huipizda007/ATB_go/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"

	// для автоматичних міграцій
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("Cannot load config: %v\n", err)
	}

	runDBMigration("file://db/migration", cfg.DBSource)

	connPool, err := pgxpool.New(context.Background(), cfg.DBSource)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer connPool.Close()

	err = connPool.Ping(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ping failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully connected to the database!")

	store := db.NewStore(connPool)
	server := server.NewServer(store)

	fmt.Printf("Starting server on %s...\n", cfg.ServerAddress)
	server.Run(cfg.ServerAddress)
}

// runDBMigration читає SQL-файли і застосовує їх до бази даних
func runDBMigration(migrationURL string, dbSource string) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		log.Fatalf("Cannot create new migrate instance: %v", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Failed to run migrate up: %v", err)
	}

	fmt.Println("DB migrated successfully!")
}