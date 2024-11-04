package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/rikiitokazu/go-backend/internal/api"
	"github.com/rikiitokazu/go-backend/internal/db"
	"github.com/rikiitokazu/go-backend/internal/logger"
	"github.com/stripe/stripe-go/v78"
)

// TODO: use log.Fatal vs log.Panic
// TODO: fix air config
// TODO: It is actually a pretty bad idea in microservices to "just exit" when encountering an error such as an external resource being unavailable available (database, cache, ...). The reason for this is that when this external resource suddenly causes errors, it is very likely that this will be an issue for many other microservice instances, causing a storm of microservices restarting, making troubleshooting issues and recovering much harder. What you ideally need to do is internally retry with a backoff timer and reconnect yourself.*/
// NOTE: init function runs before main

/**
 * Main function
 */
func main() {
	logger.InitLogger()
	defer logger.CloseLogger()

	// load env variables
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load env variables", err)
	}

	// Initalize database
	db := db.CreateDatabase()

	// TODO: defer db.Close()

	// Initialize repositories (data-access layer)
	repo := api.InitializeRepositories(db)

	// Initialize handlers (handler layer)
	handlers := api.InitializeHandlers(repo)

	// Initializaing app
	app := api.CreateNewApp(db, handlers)

	// Loads Stripe
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Handles context signals/interruptions
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Starts App
	err = app.Start(ctx)
	if err != nil {
		log.Fatal("failed to start app", err)
	}

}
