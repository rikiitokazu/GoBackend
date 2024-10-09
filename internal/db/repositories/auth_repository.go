package repositories

import (
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/rikiitokazu/go-backend/internal/api/models"
)

type AuthRepositoryInterface interface {
	Create(*models.User) error
}

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(db *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{
		db: db,
	}
}

func (ar *AuthRepository) Create(user *models.User) error {
	return nil
}
