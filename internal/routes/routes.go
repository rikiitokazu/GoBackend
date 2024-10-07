package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rikiitokazu/go-backend/internal/validation"
)

func (a *App) loadRoutes() {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// allow CORS
	router.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// Base Router
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// Subrouters
	// router.Route("/config", loadConfigRoute)
	// router.Route("/create-checkout-session", a.loadCheckoutRoute)
	// router.Route("/session-status", loadRetrieveRoute)
	router.Route("/user_profile", a.loadUserSetupRoutes)
	// router.Route("/webhook", a.loadWebhookRouter)

	/*
	 Ensuring we can use the router across functions,
	 as long as it is a function with a receiver argument
	*/
	a.router = router
}

/*
Retruns stripe publishable key

	Could make it public on frontend if you wanted
*/
// func loadConfigRoute(router chi.Router) {
// 	router.Get("/", handler.HandleConfig)
// }

// // Loads checkout router for payment
// func (a *App) loadCheckoutRoute(router chi.Router) {
// 	paymentHandler := &handler.Payment{
// 		DatabaseConn: &database.Database{
// 			Pool: a.db,
// 		},
// 	}

// 	router.Post("/", paymentHandler.CreateCheckoutSession)
// }

// // Handles router after payment has been made
// func loadRetrieveRoute(router chi.Router) {
// 	router.Get("/", handler.RetrieveCheckoutSession)
// }

// // All functions/routes related to the database
func (a *App) loadUserSetupRoutes(router chi.Router) {

	// router.Post("/check-active-user", databaseConn.CheckActiveUser)
	// router.Post("/verify-email", databaseConn.VerifyEmail)
	router.Post("/register", validation.RegisterUser)
	router.Post("/login", validation.Login)

	// router.Post("/user_course", databaseConn.GetUserCourse)
}

// /*
// Handles webhooks, such as when a payment is successful to

// 	create a new customer
// */
// func (a *App) loadWebhookRouter(router chi.Router) {
// 	webhookHandler := &handler.Payment{
// 		DatabaseConn: &database.Database{
// 			Pool: a.db,
// 		},
// 	}
// 	router.Post("/", webhookHandler.HandleWebhook)
// }