package course

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rikiitokazu/go-backend/internal/api/models"
)

func (ch *CourseHandler) DropCourse(w http.ResponseWriter, r *http.Request) {
	// Validate Request
	var req models.CourseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	req.Date = time.Now()

	// TODO: Valid jwt function --> move to separate util
	cookie, err := r.Cookie("Authorization")
	if err != nil {
		log.Println("Couldn't receive cookie")
		return
	}
	tokenString := cookie.Value
	token, err := verifyToken(tokenString)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		log.Println("Invalid jwt")
		return
	}
	log.Println("Valid jwt")

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		log.Println("Couldn't get claims")
		return
	}
	currentUserId := claims["id"].(float64)

	// Call Drop, check error
	err = ch.CourseRepository.DropCourse(&req, currentUserId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return response
	response := struct {
		Status        string `json:"status"`
		CourseDropped int    `json:"course_dropped"`
	}{
		Status:        "success",
		CourseDropped: req.CourseNumber,
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
