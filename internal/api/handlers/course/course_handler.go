package course

import (
	"github.com/rikiitokazu/go-backend/internal/db/repositories"
)

type CourseHandler struct {
	CourseRepository *repositories.CourseRepository
}

func NewCourseHandler(cr *repositories.CourseRepository) *CourseHandler {
	return &CourseHandler{
		CourseRepository: cr,
	}
}
