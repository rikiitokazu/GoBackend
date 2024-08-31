package application

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type App struct {
	router http.Handler
	db     *pgxpool.Pool
}

func initializeDatabase() *pgxpool.Pool {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
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

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("failed to make connection pool")
	}

	fmt.Println("Database connected")
	// defer pool.Close()

	return pool
}

func New() *App {
	app := &App{
		db: initializeDatabase(),
	}
	//app.loadRoutes()
	return app
}

func (a *App) Start(ctx context.Context) error {
	server := &http.Server{
		Addr:    "localhost:8000",
		Handler: a.router,
	}

	fmt.Println("Starting Server...")
	ch := make(chan error, 1)

	go func() {
		err := server.ListenAndServe()
		if err != nil {
			ch <- fmt.Errorf("failed: %w", err)
		}
		close(ch)
	}()

	select {
	case err := <-ch:
		return err
	case <-ctx.Done():
		timeout, cancel := context.WithTimeout(context.Background(), time.Second*10)
		defer cancel()
		return server.Shutdown(timeout)
	}

}

// query := `
// CREATE TABLE IF NOT EXISTS users (
// 	id SERIAL PRIMARY KEY,
// 	name TEXT NOT NULL,
// 	email TEXT NOT NULL,
// 	password TEXT NOT NULL,
// 	registered_courses TEXT[],
// 	date_created TIMESTAMP NOT NULL,
// 	last_active TIMESTAMP NOT NULL
// )`

// _, err := pool.Exec(context.Background(), query)
// if err != nil {
// 	log.Printf(err.Error())
// 	return
// }
