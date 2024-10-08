package auth

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rikiitokazu/go-backend/models"
	"github.com/rikiitokazu/go-backend/models/user_profile_db"
)

type Database struct {
	Pool *pgxpool.Pool
}

// TODO: *time.Time or time.Time
// TODO: On the frontend, strip any white spaces
// TODO: Don't use BadRequest all the time, use what is appropriate
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	//pool := application.DB
	var req models.User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Check if the request has all the valid information (should be checked on frontend)
	if req.Password == "" || req.Email == "" || req.Name == "" {
		http.Error(w, "Fill in all the information", http.StatusBadRequest)
		return
	}
	response := struct {
		UserInfo models.User `json:"user"`
		Status   string      `json:"status"`
	}{
		UserInfo: req,
		Status:   "success",
	}
	successStatus := user_profile_db.RegisterUserInData(&req)
	if successStatus["successStatus"] != "true" {
		http.Error(w, successStatus["successStatus"], http.StatusBadRequest)
		return
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
