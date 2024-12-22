package handlers

import (
	"github.com/rikiitokazu/go-backend/internal/api/handlers/course"
	"github.com/rikiitokazu/go-backend/internal/api/handlers/user"
	"github.com/rikiitokazu/go-backend/internal/api/handlers/utils"
	"github.com/rikiitokazu/go-backend/internal/db/repositories"
)

type Handlers struct {
	UserHandler   *user.UserHandler
	CourseHandler *course.CourseHandler
	UtilsHandler  *utils.UtilsHandler
}

func NewHandlers(ur *repositories.UserRepository, cr *repositories.CourseRepository) *Handlers {
	return &Handlers{
		UserHandler:   user.NewUserHandler(ur),
		CourseHandler: course.NewCourseHandler(cr),
		UtilsHandler:  utils.NewUtilsHandler(),
	}
}
