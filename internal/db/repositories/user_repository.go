package repositories

import (
	"context"
	"errors"
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

// TODO: Does models.LoginAuth need to be a pointer?
// TODO: refactor db to be a method receiver
func (ur *UserRepository) VerifyUserExists(user *models.User) error {
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
		return err
	}
	if user.Email == "" {
		return errors.New("email doesn't exist")
	}
	return nil
}

func (ur *UserRepository) Login(user *models.User) error {
	pool := ur.db
	var encryptedPassword string
	query := `SELECT id, password, registered_courses from users
	WHERE email = $1`
	err := pool.QueryRow(context.Background(), query, user.Email).
		Scan(&user.CustomerID, &encryptedPassword, &user.RegisteredCourses)
	if err != nil {
		return err
	}

	// Does the password match?
	err = bcrypt.CompareHashAndPassword([]byte(encryptedPassword), []byte(user.Password))
	if err != nil {
		return err
	}
	return nil
}

func (ur *UserRepository) Register(user *models.User) error {
	pool := ur.db

	// Determine if email already exists
	var emailExists bool
	err := pool.QueryRow(context.Background(), `SELECT EXISTS (
			SELECT 1
			FROM users
			WHERE email = $1
		)`, user.Email).Scan(&emailExists)
	if err != nil {
		return err
	}
	if emailExists {
		return errors.New("email already exists")
	}

	// use bcrypt to encrypt password
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		return err
	}

	insertQuery := `
		INSERT INTO users (name, email, password, registered_courses, date_created, last_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		`
	_, err = pool.Exec(context.Background(), insertQuery,
		user.Name, user.Email, string(passwordHash), user.RegisteredCourses, time.Now(), time.Now())
	if err != nil {
		return err
	}
	return nil
}
