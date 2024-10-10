package user

import (
	"github.com/rikiitokazu/go-backend/internal/db/repositories"
)

type UserHandler struct {
	UserRepository *repositories.UserRepository
}

func NewUserHandler(ur *repositories.UserRepository) *UserHandler {
	return &UserHandler{
		UserRepository: ur,
	}
}
