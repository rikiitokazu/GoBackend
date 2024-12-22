package utils

// ============================================
// We have a request that gets the publishable key,
// which is why we have a utils handler as well
// ============================================

import (
	"net/http"
	"os"
)

type UtilsHandler struct {
}

func NewUtilsHandler() *UtilsHandler {
	return &UtilsHandler{}
}

func (uh *UtilsHandler) GetPublishableKey(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, struct {
		PublishableKey string `json:"publishableKey"`
	}{
		PublishableKey: os.Getenv("STRIPE_PUBLISHABLE_KEY"),
	})
}
