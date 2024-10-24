package handlers

import (
	"github.com/rikiitokazu/go-backend/internal/api/handlers/course"
	"github.com/rikiitokazu/go-backend/internal/api/handlers/user"
	"github.com/rikiitokazu/go-backend/internal/db/repositories"
)

type Handlers struct {
	UserHandler   *user.UserHandler
	CourseHandler *course.CourseHandler
}

func NewHandlers(ur *repositories.UserRepository, cr *repositories.CourseRepository) *Handlers {
	return &Handlers{
		UserHandler:   user.NewUserHandler(ur),
		CourseHandler: course.NewCourseHandler(cr),
	}
}
