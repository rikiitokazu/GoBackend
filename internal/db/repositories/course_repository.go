package repositories

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rikiitokazu/go-backend/internal/api/models"
)

type CourseRepositoryInterface interface {
	Enroll(*models.EnrollRequest) error
	DropCourse(*models.EnrollRequest) error
}

type CourseRepository struct {
	db *pgxpool.Pool
}

func NewCourseRepository(db *pgxpool.Pool) *CourseRepository {
	return &CourseRepository{
		db: db,
	}
}

func (cr *CourseRepository) EnrollCourse(course *models.EnrollRequest) error {
	pool := cr.db
	// Check if the course is still available.

	// Check if course.Number is valid
	if course.CourseNumber <= 0 {
		return errors.New("invalid course number")
	}

	// TODO: Waitlist
	var students int
	var capacity int
	query := `
		SELECT students, capacity
		FROM courses
		WHERE course_number = $1
		AND active = true
	`
	err := pool.QueryRow(context.Background(), query, course.CourseNumber).Scan(&students, &capacity)
	log.Println("students", students)
	log.Println("capacity", capacity)
	if err != nil {
		return err
	}
	if students > capacity {
		return errors.New("course is full")
	}

	// Get jwt of the corresponding user, and then enroll them in a course by appending to the courses array

	// Update course capacity += 1
	return nil
}

func (cr *CourseRepository) DropCourse(*models.EnrollRequest) error {
	return nil
}
