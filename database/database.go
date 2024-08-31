package database

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

var DB *pgxpool.Pool

func InitializeDatabase() *pgxpool.Pool {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
		os.Exit(1)
	}
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("db url empty")
	}
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		log.Fatal("failed to parse dbUrl")
	}

	// Make sure you don't do := because this will not refer to the global DB var
	DB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("failed to make connection pool")
	}

	fmt.Println("Database connected")
	// defer pool.Close()

	return DB
}
