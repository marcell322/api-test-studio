package persistence

import (
	"gorm.io/gorm"

	"github.com/marcell322/api-test-studio/internal/domain/models"
	"github.com/marcell322/api-test-studio/internal/domain/repository"
)

type GormHistoryRepository struct{ db *gorm.DB }

func NewGormHistoryRepository(db *gorm.DB) repository.HistoryRepository {
	return &GormHistoryRepository{db: db}
}

func (r *GormHistoryRepository) Create(h *models.RequestHistory) error {
	return r.db.Create(h).Error
}

func (r *GormHistoryRepository) GetByID(id uint) (*models.RequestHistory, error) {
	var h models.RequestHistory
	if err := r.db.First(&h, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &h, nil
}

func (r *GormHistoryRepository) ListByUser(userID uint) ([]models.RequestHistory, error) {
	var hs []models.RequestHistory
	if err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&hs).Error; err != nil {
		return nil, err
	}
	return hs, nil
}

func (r *GormHistoryRepository) Delete(id uint) error {
	return r.db.Delete(&models.RequestHistory{}, id).Error
}