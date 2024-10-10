package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rikiitokazu/go-backend/internal/api/handlers"
	"github.com/rikiitokazu/go-backend/internal/api/routes"
	"github.com/rikiitokazu/go-backend/internal/db/repositories"
)

type App struct {
	router http.Handler
	// add other configurations here
	// TODO: add another config folder with all configurations
}

func InitializeHandlers(r *repositories.Repositories) *handlers.Handlers {
	return handlers.NewHandlers(r.UserRepository)
}

func InitializeRepositories(db *pgxpool.Pool) *repositories.Repositories {
	return repositories.NewRepositories(db)
}

func CreateNewApp(db *pgxpool.Pool, hd *handlers.Handlers) *App {
	app := &App{}
	router := routes.LoadRoutes(db, hd)
	app.router = router
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
