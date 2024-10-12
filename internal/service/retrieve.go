package service

import (
	"net/http"

	"github.com/stripe/stripe-go/v78/checkout/session"
)

func RetrieveCheckoutSession(w http.ResponseWriter, r *http.Request) {
	s, err := session.Get(r.URL.Query().Get("session_id"), nil)

	if err != nil {
		http.Error(w, "Failed to retrieve session: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if s == nil {
		http.Error(w, "Session not found", http.StatusNotFound)
		return
	}

	writeJSON(w, struct {
		Status        string `json:"status"`
		CustomerEmail string `json:"customer_email"`
	}{
		Status:        string(s.Status),
		CustomerEmail: string(s.CustomerDetails.Email),
	})

}
