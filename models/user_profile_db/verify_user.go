package user_profile_db

import (
	"context"
	"log"

	"github.com/rikiitokazu/go-backend/database"
	"github.com/rikiitokazu/go-backend/models"
	"golang.org/x/crypto/bcrypt"
)

type ErrorStatement struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Error string `json:"error"`
}

// TODO: Does models.LoginAuth need to be a pointer?
// TODO: refactor db to be a method receiver
func VerifyUserLogin(req models.LoginAuth) ErrorStatement {
	pool := database.DB
	var userID int
	var userEmail string
	query := `
		SELECT id, email
		FROM users
		WHERE email = $1
	`
	err := pool.QueryRow(context.Background(), query, req.Email).Scan(&userID, &userEmail)
	if err != nil {
		log.Println(err.Error())
		return ErrorStatement{
			Error: err.Error(),
		}
	}
	if userEmail == "" {
		return ErrorStatement{
			Error: "Email doesn't exist",
		}
	}

	// Compare the request password with corresponding password hash
	var dbPass string
	query = `SELECT password from users
	WHERE email = $1`
	err = pool.QueryRow(context.Background(), query, req.Email).Scan(&dbPass)
	if err != nil {
		return ErrorStatement{
			Error: err.Error(),
		}
	}

	// Does the password match?
	err = bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(req.Password))
	if err != nil {
		return ErrorStatement{
			Error: err.Error(),
		}
	}
	return ErrorStatement{
		Id:    userID,
		Email: userEmail,
		Error: "nil",
	}
}
