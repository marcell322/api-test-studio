package repository

import "github.com/marcell322/api-test-studio/internal/domain/models"

// SavedRequestRepository defines persistence operations for saved requests
type SavedRequestRepository interface {
	Create(r *models.SavedRequest) error
	GetByID(id uint) (*models.SavedRequest, error)
	ListByUser(userID uint) ([]models.SavedRequest, error)
	ListByCollection(collectionID uint) ([]models.SavedRequest, error)
	Update(r *models.SavedRequest) error
	Delete(id uint) error
}