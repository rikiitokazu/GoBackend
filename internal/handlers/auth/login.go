package auth 

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rikiitokazu/go-backend/models"
	"github.com/rikiitokazu/go-backend/models/user_profile_db"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginAuth
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Because we used a struct value with type string, we compare the actual "nil"
	// TODO: Is there a better way to do this?
	response := user_profile_db.VerifyUserLogin(req)
	if response.Error != "nil" {
		http.Error(w, response.Error, http.StatusBadRequest)
	}

	// Generate a jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": response.Id,
		"exp": time.Now().Add(time.Minute * 1).Unix(),
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
