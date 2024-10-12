package service

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/customer"
)

type CheckoutSessionRequest struct {
	Email        string `json:"email"`
	Name         string `json:"name"`
	CourseNumber string `json:"course_number"`
}

func CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	err := godotenv.Load()
	if err != nil {
		log.Println("Couldn't load environment vars")
		os.Exit(1)
	}
	// *Reminder to allow access to multiple langugages (japanese)
	var req CheckoutSessionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Trimming spaces from email request
	req.Email = strings.TrimSpace(req.Email)

	// Getting the correct price tag according to user selection
	priceKey := getPrice(req.CourseNumber)
	if priceKey == "error" {
		log.Println("Couldn't get accesse to price tag")
		http.Error(w, "Internal Server Error: Price tag inaccessible", http.StatusInternalServerError)
		return
	}

	// Frontend Route
	domain := os.Getenv("FRONTEND_ROUTE")

	// Creating new customer session
	customerParams := &stripe.CustomerParams{
		Name:  stripe.String(req.Name),
		Email: stripe.String(req.Email),
	}
	customerParams.AddMetadata("course_id", req.CourseNumber)
	customerResult, err := customer.New(customerParams)
	if err != nil {
		log.Printf("session.New: %v", err)
	}

	// Parameters for checkout session for customer above
	params := &stripe.CheckoutSessionParams{
		UIMode:    stripe.String("embedded"),
		ReturnURL: stripe.String(domain + "/return?session_id={CHECKOUT_SESSION_ID}"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceKey),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:     stripe.String(string(stripe.CheckoutSessionModePayment)),
		Customer: stripe.String(customerResult.ID),
	}

	// Creates a new checkout session
	s, err := session.New(params)
	if err != nil {
		log.Printf("session.New: %v", err)
	}

	// Returning the clinet secret
	if s.ClientSecret == "" {
		log.Println("ClientSecret is empty. Unable to process the payment.")
		http.Error(w, "Internal Server Error: Client secret is missing", http.StatusInternalServerError)
		return
	}
	writeJSON(w, struct {
		ClientSecret string `json:"clientSecret"`
	}{
		ClientSecret: s.ClientSecret,
	})

}
