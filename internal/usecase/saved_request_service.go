package usecase

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/marcell322/api-test-studio/internal/domain/models"
	"github.com/marcell322/api-test-studio/internal/domain/repository"
)

var validHTTPMethods = map[string]bool{
	"GET": true, "POST": true, "PUT": true, "PATCH": true, "DELETE": true,
}

type SavedRequestService interface {
	Create(userID, collectionID uint, name, method, url string, headers map[string]string, body string) (*models.SavedRequest, error)
	List(userID uint) ([]models.SavedRequest, error)
	ListByCollection(userID, collectionID uint) ([]models.SavedRequest, error)
	Get(userID, id uint) (*models.SavedRequest, error)
	Update(userID, id uint, name, method, url string, headers map[string]string, body string) (*models.SavedRequest, error)
	Delete(userID, id uint) error
}

type savedRequestService struct {
	repo           repository.SavedRequestRepository
	collectionRepo repository.CollectionRepository
}

func NewSavedRequestService(r repository.SavedRequestRepository, cr repository.CollectionRepository) SavedRequestService {
	return &savedRequestService{repo: r, collectionRepo: cr}
}

func validateSavedRequest(name, method, url string) error {
	if strings.TrimSpace(name) == "" {
		return errors.New("name is required")
	}
	if !validHTTPMethods[strings.ToUpper(strings.TrimSpace(method))] {
		return errors.New("invalid http method")
	}
	if strings.TrimSpace(url) == "" {
		return errors.New("url is required")
	}
	return nil
}

// checkCollectionOwnership ensures collectionID exists and belongs to userID.
// A request can only be created in / attached to a collection you own.
func (s *savedRequestService) checkCollectionOwnership(userID, collectionID uint) error {
	col, err := s.collectionRepo.GetByID(collectionID)
	if err != nil {
		return err
	}
	if col == nil {
		return ErrNotFound
	}
	if col.UserID != userID {
		return ErrForbidden
	}
	return nil
}

func marshalHeaders(headers map[string]string) (string, error) {
	if headers == nil {
		return "{}", nil
	}
	b, err := json.Marshal(headers)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func (s *savedRequestService) Create(userID, collectionID uint, name, method, url string, headers map[string]string, body string) (*models.SavedRequest, error) {
	if err := validateSavedRequest(name, method, url); err != nil {
		return nil, err
	}
	if err := s.checkCollectionOwnership(userID, collectionID); err != nil {
		return nil, err
	}
	headersJSON, err := marshalHeaders(headers)
	if err != nil {
		return nil, err
	}
	sr := &models.SavedRequest{
		UserID:       userID,
		CollectionID: collectionID,
		Name:         strings.TrimSpace(name),
		Method:       strings.ToUpper(strings.TrimSpace(method)),
		URL:          strings.TrimSpace(url),
		Headers:      headersJSON,
		Body:         body,
	}
	if err := s.repo.Create(sr); err != nil {
		return nil, err
	}
	return sr, nil
}

func (s *savedRequestService) List(userID uint) ([]models.SavedRequest, error) {
	return s.repo.ListByUser(userID)
}

func (s *savedRequestService) ListByCollection(userID, collectionID uint) ([]models.SavedRequest, error) {
	if err := s.checkCollectionOwnership(userID, collectionID); err != nil {
		return nil, err
	}
	return s.repo.ListByCollection(collectionID)
}

// Get enforces that the saved request belongs to userID.
func (s *savedRequestService) Get(userID, id uint) (*models.SavedRequest, error) {
	sr, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if sr == nil {
		return nil, ErrNotFound
	}
	if sr.UserID != userID {
		return nil, ErrForbidden
	}
	return sr, nil
}

func (s *savedRequestService) Update(userID, id uint, name, method, url string, headers map[string]string, body string) (*models.SavedRequest, error) {
	if err := validateSavedRequest(name, method, url); err != nil {
		return nil, err
	}
	sr, err := s.Get(userID, id)
	if err != nil {
		return nil, err
	}
	headersJSON, err := marshalHeaders(headers)
	if err != nil {
		return nil, err
	}
	sr.Name = strings.TrimSpace(name)
	sr.Method = strings.ToUpper(strings.TrimSpace(method))
	sr.URL = strings.TrimSpace(url)
	sr.Headers = headersJSON
	sr.Body = body
	if err := s.repo.Update(sr); err != nil {
		return nil, err
	}
	return sr, nil
}

func (s *savedRequestService) Delete(userID, id uint) error {
	sr, err := s.Get(userID, id)
	if err != nil {
		return err
	}
	return s.repo.Delete(sr.ID)
}