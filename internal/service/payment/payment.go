package payment

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/rikiitokazu/go-backend/internal/api/models"
	"github.com/rikiitokazu/go-backend/internal/service"
	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/checkout/session"
	"github.com/stripe/stripe-go/v78/customer"
)

// We do not need db access for stripe

func CreateCheckoutSession(courseNum int, userInfo *models.User) error {
	// TODO: Reminder to allow access to multiple langugages (japanese)
	// TODO: Yen
	domain := os.Getenv("FRONTEND_ROUTE")

	// 1) Get the correct course based on courseNum
	// 2) Register that user for the course
	// 2.1) if not using course db, make sure that we check if its full first, but we did that in db step
	// Getting the correct price tag according to user selection
	// 3) Add that user as a customer that is associated with the course purchase

	priceKey := service.GetPrice(courseNum)
	if priceKey == "error" {
		log.Println("Couldn't get accesse to price tag")
		// http.Error(w, "Internal Server Error: Price tag inaccessible", http.StatusInternalServerError)
		return nil
	}

	// Creating new customer session
	customerParams := &stripe.CustomerParams{
		Name:  stripe.String(userInfo.Name),
		Email: stripe.String(userInfo.Email),
	}
	courseToString := strconv.Itoa(courseNum)
	customerParams.AddMetadata("course_id", courseToString)
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

	// Returning the client secret
	if s.ClientSecret == "" {
		log.Println("ClientSecret is empty. Unable to process the payment.")
		return errors.New("client secret is empty")
	}
	return nil

}
