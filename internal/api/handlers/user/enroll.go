package user

import (
	"encoding/json"
	"net/http"

	"github.com/rikiitokazu/go-backend/internal/api/models"
)

func (uh *UserHandler) EnrollCourse(w http.ResponseWriter, r *http.Request) {
	var req models.EnrollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check availability of course

	// Enroll in stripe, if it is not free

	// Add to database

	// Return http response
}
