package repositories

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	UserRepository   *UserRepository
	CourseRepository *CourseRepository
}

// TODO: SERIAL or UUID
func NewRepositories(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		UserRepository:   NewUserRepository(db),
		CourseRepository: NewCourseRepository(db),
	}
}
