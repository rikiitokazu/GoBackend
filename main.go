package main

import (
	"context"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/rikiitokazu/go-backend/application"
	"github.com/rikiitokazu/go-backend/database"
	"github.com/stripe/stripe-go/v78"
)

// TODO: use log.Fatal vs log.Panic
// TODO: error with localhost:8000 not gracefully shutdown
/* TODO: It is actually a pretty bad idea in microservices to "just exit" when encountering an error such as an external resource being unavailable available (database, cache, ...). The reason for this is that when this external resource suddenly causes errors, it is very likely that this will be an issue for many other microservice instances, causing a storm of microservices restarting, making troubleshooting issues and recovering much harder. What you ideally need to do is internally retry with a backoff timer and reconnect yourself.*/
// TODO: 'air' doesn't work, localhost is not terminated on graceful shutdown

/**
 * Initial function that runs before main to load env vars
 */
func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("failed to load env variables")
	}
}

/**
 * Main function
 */
func main() {
	// Initializaing app
	app := application.New()

	// Initalize database
	database.InitializeDatabase()

	// Loads Stripe
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// Handles context signals/interruptions
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// Starts App
	err := app.Start(ctx)
	if err != nil {
		log.Fatal("failed to start app", err)
	}

}
