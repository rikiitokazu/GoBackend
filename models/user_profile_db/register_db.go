package user_profile_db

import (
	"context"
	"log"
	"time"

	"github.com/rikiitokazu/go-backend/database"
	"github.com/rikiitokazu/go-backend/models"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUserInData(req *models.User) map[string]string {
	if database.DB == nil {
		log.Println("FAILURE")
		var maps map[string]string
		return maps
	}
	pool := database.DB
	userSuccess := make(map[string]string)
	// Determine if email already exists
	var emailExists bool
	err := pool.QueryRow(context.Background(), `SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE email = $1
		)`, req.Email).Scan(&emailExists)
	if err != nil {
		//http.Error(w, err.Error(), http.StatusBadRequest)
		userSuccess["successStatus"] = err.Error()
		return userSuccess
	}
	if emailExists {
		userSuccess["successStatus"] = "Email already exists"
		return userSuccess
	}

	// use bcrypt to encrypt password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), 10)
	if err != nil {
		userSuccess["successStatus"] = err.Error()
		return userSuccess
	}

	insertQuery := `
		INSERT INTO users (name, email, password, registered_courses, date_created, last_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		`
	_, err = pool.Exec(context.Background(), insertQuery, req.Name, req.Email, string(passwordHash), req.RegisteredCourses, time.Now(), time.Now())
	if err != nil {
		userSuccess["successStatus"] = "Error occured while inserting data"
		return userSuccess
	}
	userSuccess["successStatus"] = "true"
	return userSuccess
}
