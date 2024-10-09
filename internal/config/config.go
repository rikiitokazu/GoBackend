package config

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rikiitokazu/go-backend/internal/routes"
)

type App struct {
	router http.Handler
}

func CreateNewApp() *App {
	app := &App{}
	router := routes.LoadRoutes()
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
