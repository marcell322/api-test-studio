package repository

import "github.com/marcell322/api-test-studio/internal/domain/models"

// CollectionRepository defines persistence operations for collections
type CollectionRepository interface {
	Create(c *models.Collection) error
	GetByID(id uint) (*models.Collection, error)
	ListByUser(userID uint) ([]models.Collection, error)
	Update(c *models.Collection) error
	Delete(id uint) error
}