package handlers

import (
	"github.com/rikiitokazu/go-backend/internal/db/repositories"
)

type Handlers struct {
	AuthHandler *AuthHandler
}

func NewHandlers(ar *repositories.AuthRepository) *Handlers {
	return &Handlers{
		AuthHandler: NewAuthHandler(ar),
	}
}
