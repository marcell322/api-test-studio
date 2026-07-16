package persistence

import (
	"gorm.io/gorm"

	"github.com/marcell322/api-test-studio/internal/domain/models"
	"github.com/marcell322/api-test-studio/internal/domain/repository"
)

type GormUserRepository struct{ db *gorm.DB }

func NewGormUserRepository(db *gorm.DB) repository.UserRepository { return &GormUserRepository{db: db} }

func (r *GormUserRepository) CreateUser(u *models.User) error { return r.db.Create(u).Error }

func (r *GormUserRepository) GetUserByEmail(email string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		if err == gorm.ErrRecordNotFound { return nil, nil }
		return nil, err
	}
	return &u, nil
}

func (r *GormUserRepository) GetUserByID(id uint) (*models.User, error) {
	var u models.User
	if err := r.db.First(&u, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound { return nil, nil }
		return nil, err
	}
	return &u, nil
}

func (r *GormUserRepository) CheckEmailExists(email string) (bool, error) {
	var count int64
	if err := r.db.Model(&models.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
