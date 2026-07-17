package persistence

import (
	"gorm.io/gorm"

	"github.com/marcell322/api-test-studio/internal/domain/models"
	"github.com/marcell322/api-test-studio/internal/domain/repository"
)

type GormSavedRequestRepository struct{ db *gorm.DB }

func NewGormSavedRequestRepository(db *gorm.DB) repository.SavedRequestRepository {
	return &GormSavedRequestRepository{db: db}
}

func (r *GormSavedRequestRepository) Create(sr *models.SavedRequest) error {
	return r.db.Create(sr).Error
}

func (r *GormSavedRequestRepository) GetByID(id uint) (*models.SavedRequest, error) {
	var sr models.SavedRequest
	if err := r.db.First(&sr, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &sr, nil
}

func (r *GormSavedRequestRepository) ListByUser(userID uint) ([]models.SavedRequest, error) {
	var rs []models.SavedRequest
	if err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&rs).Error; err != nil {
		return nil, err
	}
	return rs, nil
}

func (r *GormSavedRequestRepository) ListByCollection(collectionID uint) ([]models.SavedRequest, error) {
	var rs []models.SavedRequest
	if err := r.db.Where("collection_id = ?", collectionID).Order("created_at desc").Find(&rs).Error; err != nil {
		return nil, err
	}
	return rs, nil
}

func (r *GormSavedRequestRepository) Update(sr *models.SavedRequest) error {
	return r.db.Save(sr).Error
}

func (r *GormSavedRequestRepository) Delete(id uint) error {
	return r.db.Delete(&models.SavedRequest{}, id).Error
}