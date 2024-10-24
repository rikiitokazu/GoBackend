package course

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/rikiitokazu/go-backend/internal/api/models"
)

func (ch *CourseHandler) EnrollCourse(w http.ResponseWriter, r *http.Request) {
	var req models.EnrollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check availability of course in "courses" table
	err := ch.CourseRepository.EnrollCourse(&req)
	if err != nil {
		log.Println("error")
	}
	// Enroll in stripe, if it is not free

	// Add to database

	// Return http response
}
