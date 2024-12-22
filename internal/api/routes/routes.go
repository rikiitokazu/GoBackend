package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/rikiitokazu/go-backend/internal/api/handlers"
)

func LoadRoutes(db *pgxpool.Pool, h *handlers.Handlers) *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)

	// allow CORS
	router.Use(cors.Handler(cors.Options{
		// the receiver's router field, so that it can be used in other
		// functions.
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
	router.Route("/user", func(router chi.Router) {
		loadUserSetupRoutes(router, h)
	})
	router.Route("/course", func(router chi.Router) {
		loadCourseSetupRoutes(router, h)
	})
	router.Route("/util", func(router chi.Router) {
		loadUtilSetupRoutes(router, h)
	})

	// router.Route("/metrics", promhttp.Handler())
	router.Route("/metrics", func(router chi.Router) {
		loadPrometheusRoute(router, h)
	})

	// router.Route("/webhook", a.loadWebhookRouter)

	/*
	 Ensuring we can use the router across functions,
	 as long as it is a function with a receiver argument
	*/

	return router
}

// Loads checkout router for payment
// func (a *App) loadCheckoutRoute(router chi.Router) {
// 	paymentHandler := &handler.Payment{
// 		DatabaseConn: &database.Database{
// 			Pool: a.db,
// 		},
// 	}

// 	router.Post("/", paymentHandler.CreateCheckoutSession)
// }

// Handles router after payment has been made
//
//	func loadRetrieveRoute(router chi.Router) {
//		router.Get("/", handler.RetrieveCheckoutSession)
//	}
func loadUtilSetupRoutes(router chi.Router, h *handlers.Handlers) {
	// GET -- Stripe public key
	router.Get("/payment-config", h.UtilsHandler.GetPublishableKey)
}

func loadUserSetupRoutes(router chi.Router, h *handlers.Handlers) {

	// router.Post("/check-active-user", databaseConn.CheckActiveUser)
	// router.Post("/verify-email", databaseConn.VerifyEmail)
	router.Post("/register", h.UserHandler.RegisterUser)
	router.Post("/login", h.UserHandler.LoginUser)

	// router.Post("/user_course", databaseConn.GetUserCourse)
}

func loadCourseSetupRoutes(router chi.Router, h *handlers.Handlers) {
	router.Post("/enroll", h.CourseHandler.EnrollCourse)
	router.Delete("/drop", h.CourseHandler.DropCourse)
}

func loadPrometheusRoute(router chi.Router, h *handlers.Handlers) {
	router.Handle("/", promhttp.Handler())
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
