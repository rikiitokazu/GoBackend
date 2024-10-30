package user

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rikiitokazu/go-backend/internal/api/models"
)

// TODO: do http methods need to return something?
func (uh *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Because we used a struct value with type string, we compare the actual "nil"
	// TODO: Is there a better way to do this?
	// TODO: is the return value necessary, or does the http.Error do it?
	err := uh.UserRepository.VerifyUserExists(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Given user exists, validate is password is correct
	err = uh.UserRepository.Login(&req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":    req.CustomerID,
		"name":  req.Name,
		"email": req.Email,
		"exp":   time.Now().Add(time.Minute * 1).Unix(),
	})
	tokenString, err := token.SignedString([]byte(os.Getenv("TOKEN_SECRET")))

	if err != nil {
		log.Println(err.Error())
		http.Error(w, "Failed to parse token", http.StatusBadGateway)
		return
	}
	res := struct {
		Token string `json:"token"`
	}{
		Token: tokenString,
	}
	jsonResponse, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "failed to marshal json", http.StatusBadRequest)
		return
	}
	// send it back
	// TODO: cookies or cache
	cookie := http.Cookie{
		Name:     "Authorization",
		Value:    res.Token,
		Path:     "/",
		MaxAge:   3600,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, &cookie)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
