package service

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/stripe/stripe-go/v78"
	"github.com/stripe/stripe-go/v78/price"
)

// TODO: Should all http requests implement this?
func writeJSON(w http.ResponseWriter, v interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("json.NewEncoder.Encode: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := io.Copy(w, &buf); err != nil {
		log.Printf("io.Copy: %v", err)
		return
	}
}

func GetPrice(courseNumber string) string {
	priceParams := &stripe.PriceSearchParams{
		SearchParams: stripe.SearchParams{
			Query: "active:'true'",
		},
	}

	result := price.Search(priceParams)

	for result.Next() {
		p := result.Price()
		if p.Metadata["course_id"] == courseNumber {
			return p.ID
		}
	}
	return "error"

}
