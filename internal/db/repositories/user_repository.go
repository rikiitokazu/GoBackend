package repositories

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"github.com/rikiitokazu/go-backend/internal/api/models"
)

type UserRepositoryInterface interface {
	Register(*models.User) error
}

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (ur *UserRepository) VerifyUserExists(user *models.User) models.ErrorStatement {
	pool := ur.db
	var userID int
	var userEmail string
	query := `
		SELECT id, email
		FROM users
		WHERE email = $1
	`
	err := pool.QueryRow(context.Background(), query, user.Email).Scan(&userID, &userEmail)
	if err != nil {
		log.Println(err.Error())
		return models.ErrorStatement{
			Error: err.Error(),
		}
	}
	if userEmail == "" {
		return models.ErrorStatement{
			Error: "Email doesn't exist",
		}
	}

	// Compare the request password with corresponding password hash
	var dbPass string
	query = `SELECT password from users
	WHERE email = $1`
	err = pool.QueryRow(context.Background(), query, user.Email).Scan(&dbPass)
	if err != nil {
		return models.ErrorStatement{
			Error: err.Error(),
		}
	}

	// Does the password match?
	err = bcrypt.CompareHashAndPassword([]byte(dbPass), []byte(user.Password))
	if err != nil {
		return models.ErrorStatement{
			Error: err.Error(),
		}
	}
	return models.ErrorStatement{
		Id:    userID,
		Email: userEmail,
		Error: "nil",
	}
}

func (ur *UserRepository) Register(user *models.User) error {
	pool := ur.db
	if pool == nil {
		log.Fatal("FAILURE")
		var maps map[string]string
		return maps
	}
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
