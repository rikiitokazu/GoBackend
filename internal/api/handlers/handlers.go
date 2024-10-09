package handlers

import (
	"github.com/rikiitokazu/go-backend/internal/api/handlers/user"
	"github.com/rikiitokazu/go-backend/internal/db/repositories"
)

type Handlers struct {
	UserHandler *user.UserHandler
}

func NewHandlers(ur *repositories.UserRepository) *Handlers {
	return &Handlers{
		UserHandler: user.NewUserHandler(ur),
	}
}
