package persistence

import (
	"gorm.io/gorm"

	"github.com/marcell322/api-test-studio/internal/domain/models"
	"github.com/marcell322/api-test-studio/internal/domain/repository"
)

type GormCollectionRepository struct{ db *gorm.DB }

func NewGormCollectionRepository(db *gorm.DB) repository.CollectionRepository {
	return &GormCollectionRepository{db: db}
}

func (r *GormCollectionRepository) Create(c *models.Collection) error {
	return r.db.Create(c).Error
}

func (r *GormCollectionRepository) GetByID(id uint) (*models.Collection, error) {
	var c models.Collection
	if err := r.db.First(&c, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *GormCollectionRepository) ListByUser(userID uint) ([]models.Collection, error) {
	var cs []models.Collection
	if err := r.db.Where("user_id = ?", userID).Order("created_at desc").Find(&cs).Error; err != nil {
		return nil, err
	}
	return cs, nil
}

func (r *GormCollectionRepository) Update(c *models.Collection) error {
	return r.db.Save(c).Error
}

func (r *GormCollectionRepository) Delete(id uint) error {
	return r.db.Delete(&models.Collection{}, id).Error
}