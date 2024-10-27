package course

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rikiitokazu/go-backend/internal/api/models"
)

func (ch *CourseHandler) EnrollCourse(w http.ResponseWriter, r *http.Request) {
	var req models.EnrollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Valid jwt
	cookie, err := r.Cookie("Authorization")
	if err != nil {
		log.Println("Couldn't receive cookie")
		return
	}
	tokenString := cookie.Value
	err = verifyToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("Invalid jwt")
		return
	}
	log.Println("Valid jwt")
	// Check availability of course in "courses" table
	err = ch.CourseRepository.EnrollCourse(&req)
	if err != nil {
		log.Println("error")
	}
	// Enroll in stripe, if it is not free

	// Add to database

	// Return http response
}

// TODO: move to utils
func verifyToken(tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("TOKEN_SECRET")), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("invalid jwt")
	}

	return nil
}
