package repository

import "github.com/marcell322/api-test-studio/internal/domain/models"

// UserRepository defines persistence operations for users
type UserRepository interface {
	CreateUser(u *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id uint) (*models.User, error)
	CheckEmailExists(email string) (bool, error)
}
