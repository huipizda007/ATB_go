package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	db "github.com/huipizda007/ATB_go/db/sqlc"
	"github.com/huipizda007/ATB_go/internal/server"
)

func main() {
	connStr := "postgresql://postgres:postgres@127.0.0.1:5432/fruits?sslmode=disable"
	connPool, err := pgxpool.New(context.Background(), connStr)
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
	
	fmt.Println("Starting server on :3000...")
	server.Run(":3000")
}