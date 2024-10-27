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

func (cr *CourseRepository) EnrollCourse(course *models.EnrollRequest, userId float64) error {
	pool := cr.db

	// Check if course.Number is valid
	if course.CourseNumber <= 0 {
		return errors.New("invalid course number")
	}

	// Check if the course is still available.
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
	// Enroll in stripe, if it is not free
	// TODO: For right now, lets assume only course "1" is free
	// if req.CourseNumber != 1 {
	// 	log.Println("Enrolling in stripe")
	// }

	// Check if user is already enrolled in the course
	err = cr.checkUserInCourse(userId)
	if err != nil {
		return err
	}

	// Update course capacity += 1
	err = cr.addCountToCourse(course.CourseNumber)
	if err != nil {
		return err
	}
	// Update student count += 1, and append to registered_courses in users table

	return nil
}

func (cr *CourseRepository) checkUserInCourse(userId float64) error {
	pool := cr.db

	var emailExists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE id = $1);`
	err := pool.QueryRow(context.Background(), query, userId).Scan(&emailExists)
	if err != nil {
		return err
	}
	return nil
}

func (cr *CourseRepository) addCountToCourse(courseNum int) error {
	pool := cr.db
	query := `
	UPDATE courses
	SET students = students + 1
	WHERE course_number = $1 and active = true
	`
	_, err := pool.Exec(context.Background(), query, courseNum)
	if err != nil {
		return err
	}
	return nil
}

func (cr *CourseRepository) addCountToUserArray(userId float64, courseNum int) error {
	return nil
}
func (cr *CourseRepository) DropCourse(*models.EnrollRequest) error {
	return nil
}
