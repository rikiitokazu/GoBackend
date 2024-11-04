package repositories

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rikiitokazu/go-backend/internal/api/models"
	"github.com/rikiitokazu/go-backend/internal/service/payment"
)

type CourseRepositoryInterface interface {
	Enroll(*models.CourseRequest) error
	DropCourse(*models.CourseRequest) error
}

type CourseRepository struct {
	db *pgxpool.Pool
}

func NewCourseRepository(db *pgxpool.Pool) *CourseRepository {
	return &CourseRepository{
		db: db,
	}
}

func (cr *CourseRepository) EnrollCourse(course *models.CourseRequest, userInfo *models.User) error {
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

	// Check if user is already enrolled in the course
	//TODO: customerId versus userId --> keep homogenous
	err = cr.checkUserInCourse(userInfo.CustomerID, course.CourseNumber)
	if err != nil {
		return err
	}
	// Enroll in stripe, if it is not free
	// TODO: For right now, lets assume only course "1" is free
	if course.CourseNumber != 1 {
		err := payment.CreateCheckoutSession(course.CourseNumber, userInfo)
		if err != nil {
			return err
		}
	}

	// Update course capacity += 1
	err = cr.addCountToCourse(course.CourseNumber)
	if err != nil {
		return err
	}
	// Update student count += 1, and append to registered_courses in users table
	err = cr.addCountToUserArray(userInfo.CustomerID, course.CourseNumber)
	if err != nil {
		return err
	}
	return nil
}

// TODO: userId is currently float64 as a result of when we initialized it
// TODO: return (error, HTTPstatus)
func (cr *CourseRepository) checkUserInCourse(userId float64, courseNum int) error {
	pool := cr.db

	var userExists bool
	query := `SELECT EXISTS(
	SELECT id, registered_courses
	FROM users 
	WHERE id = $1 AND $2 = ANY(registered_courses)
	);`
	err := pool.QueryRow(context.Background(), query, userId, courseNum).Scan(&userExists)
	if err != nil {
		return err
	}
	if userExists {
		return errors.New("user is already enrolled")
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
	// TODO: create column for users
	return nil
}

func (cr *CourseRepository) addCountToUserArray(userId float64, courseNum int) error {
	pool := cr.db
	query := `
	UPDATE users
	SET registered_courses = ARRAY_APPEND(registered_courses, $1)
	WHERE id = $2
	`
	_, err := pool.Exec(context.Background(), query, courseNum, userId)
	if err != nil {
		return err
	}
	return nil
}

// TODO: Instead of active column, we could just do if it falls within a timeframe
func (cr *CourseRepository) DropCourse(course *models.CourseRequest, userId float64) error {
	pool := cr.db
	// Check if the course exists
	var courseExists bool
	query := `SELECT EXISTS(
	SELECT course_number, active
	FROM courses
	WHERE course_number = $1 AND active = true
	);`
	err := pool.QueryRow(context.Background(), query, course.CourseNumber).Scan(&courseExists)
	if err != nil {
		return err
	}
	if !courseExists {
		return errors.New("course does not exist")
	}

	// TODO: Do we need to check if the user is actually enrolled, or are we assuming it
	// will always be valid?

	// Remove the element from the registered_courses column in the users table
	query = `
	UPDATE users
	SET registered_courses = ARRAY_REMOVE(registered_courses, $1)
	WHERE id = $2
	`
	_, err = pool.Exec(context.Background(), query, course.CourseNumber, userId)
	if err != nil {
		return err
	}

	// Subtract the number of enrolled people of that current, active course
	query = `
	UPDATE courses
	SET students = students - 1
	WHERE course_number = $1 AND active = true
	`
	_, err = pool.Exec(context.Background(), query, course.CourseNumber)
	if err != nil {
		return err
	}

	return nil
}
