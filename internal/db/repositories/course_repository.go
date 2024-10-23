package repositories

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rikiitokazu/go-backend/internal/api/models"
)

type CourseRepositoryInterface interface {
	Enroll(*models.EnrollRequest) error
}

type CourseRepository struct {
	db *pgxpool.Pool
}

func NewCourseRepository(db *pgxpool.Pool) *CourseRepository {
	return &CourseRepository{
		db: db,
	}
}

func (cr *CourseRepository) EnrollCourse(user *models.EnrollRequest) error {
	pool := cr.db
	var userID int
	var userEmail string
	query := `
		SELECT id, email
		FROM users
		WHERE email = $1
	`
	err := pool.QueryRow(context.Background(), query, "placeholder").Scan(&userID, &userEmail)
	if err != nil {
		return err
	}
	return nil
}
