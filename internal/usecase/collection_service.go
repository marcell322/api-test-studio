package usecase

import (
	"errors"
	"strings"

	"github.com/marcell322/api-test-studio/internal/domain/models"
	"github.com/marcell322/api-test-studio/internal/domain/repository"
)

// Shared usecase-level errors, checked by handlers to pick the right HTTP status.
var (
	ErrNotFound = errors.New("not found")
	ErrForbidden = errors.New("forbidden")
)

// CollectionService defines business use-cases around collections
type CollectionService interface {
	Create(userID uint, name string) (*models.Collection, error)
	List(userID uint) ([]models.Collection, error)
	Get(userID, id uint) (*models.Collection, error)
	Rename(userID, id uint, name string) (*models.Collection, error)
	Delete(userID, id uint) error
}

type collectionService struct {
	repo repository.CollectionRepository
}

func NewCollectionService(r repository.CollectionRepository) CollectionService {
	return &collectionService{repo: r}
}

func (s *collectionService) Create(userID uint, name string) (*models.Collection, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("name is required")
	}
	c := &models.Collection{UserID: userID, Name: name}
	if err := s.repo.Create(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *collectionService) List(userID uint) ([]models.Collection, error) {
	return s.repo.ListByUser(userID)
}

// Get fetches a collection and enforces that it belongs to userID.
func (s *collectionService) Get(userID, id uint) (*models.Collection, error) {
	c, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, ErrNotFound
	}
	if c.UserID != userID {
		return nil, ErrForbidden
	}
	return c, nil
}

func (s *collectionService) Rename(userID, id uint, name string) (*models.Collection, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("name is required")
	}
	c, err := s.Get(userID, id)
	if err != nil {
		return nil, err
	}
	c.Name = name
	if err := s.repo.Update(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *collectionService) Delete(userID, id uint) error {
	c, err := s.Get(userID, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(c.ID)
}