package repository

import "github.com/marcell322/api-test-studio/internal/domain/models"

// HistoryRepository defines persistence operations for request history
type HistoryRepository interface {
	Create(h *models.RequestHistory) error
	GetByID(id uint) (*models.RequestHistory, error)
	ListByUser(userID uint) ([]models.RequestHistory, error)
	Delete(id uint) error
}