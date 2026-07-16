package usecase

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/marcell322/api-test-studio/internal/domain/models"
	"github.com/marcell322/api-test-studio/internal/domain/repository"
)

// UserService defines business use-cases around users
type UserService interface {
	Register(username, email, password string) (*models.User, error)
	Authenticate(email, password string) (*models.User, error)
}

type userService struct{ repo repository.UserRepository }

func NewUserService(r repository.UserRepository) UserService { return &userService{repo: r} }

func (s *userService) Register(username, email, password string) (*models.User, error) {
	if username == "" || email == "" || password == "" {
		return nil, errors.New("missing fields")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil { return nil, err }
	u := &models.User{Username: username, Email: email, PasswordHash: string(hash)}
	if err := s.repo.CreateUser(u); err != nil { return nil, err }
	return u, nil
}

func (s *userService) Authenticate(email, password string) (*models.User, error) {
	u, err := s.repo.GetUserByEmail(email)
	if err != nil { return nil, err }
	if u == nil { return nil, errors.New("invalid credentials") }
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid credentials")
	}
	return u, nil
}
