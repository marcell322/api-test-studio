package repository

import "github.com/marcell322/api-test-studio/internal/domain/models"

// UserRepository defines persistence operations for users
type UserRepository interface {
	Create(u *models.User) error
	GetByEmail(email string) (*models.User, error)
	GetByID(id uint) (*models.User, error)
}
