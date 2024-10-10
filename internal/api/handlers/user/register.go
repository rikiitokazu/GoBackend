package user

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/rikiitokazu/go-backend/internal/api/models"
)

// TODO: *time.Time or time.Time
// TODO: On the frontend, strip any white spaces
// TODO: Don't use BadRequest all the time, use what is appropriate
func (uh *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Check if the request has all the valid information (should be checked on frontend)
	if user.Password == "" || user.Email == "" || user.Name == "" {
		http.Error(w, "Fill in all the information", http.StatusBadRequest)
		return
	}

	// TODO: bcrypt to encrypt password
	// passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }
	// req.Password = string(passwordHash)
	// user := models.User{
	// 	Name:     req.Name,
	// 	Email:    req.Email,
	// 	Password: req.Password,
	// }
	// user.Password= newPassword thats encrypted

	if err := uh.UserRepository.Register(&user); err != nil {
		log.Fatal("failed to create user in db", err)
		return
	}

	response := struct {
		UserInfo models.User `json:"user"`
		Status   string      `json:"status"`
	}{
		UserInfo: user,
		Status:   "success",
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)

	// TODO: return value?
}
