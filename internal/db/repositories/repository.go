package repositories

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repositories struct {
	AuthRepository *AuthRepository
	// CourseRepository *CourseRepository
}

func NewRepositories(db *pgxpool.Pool) *Repositories {
	return &Repositories{
		AuthRepository: NewAuthRepository(db),
		// CourseRepository: NewCourseRepository(db),
	}
}
