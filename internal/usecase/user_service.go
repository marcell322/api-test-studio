package usecase

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

	"github.com/marcell322/api-test-studio/internal/adapters/auth"
	"github.com/marcell322/api-test-studio/internal/domain/models"
	"github.com/marcell322/api-test-studio/internal/domain/repository"
)

// UserService defines business use-cases around users
type UserService interface {
	Register(username, email, password string) (*models.User, error)
	Login(email, password string) (string, *models.User, error)
}

type userService struct{
	repo      repository.UserRepository
	jwtSecret string
	jwtHours  int
}

func NewUserService(r repository.UserRepository, jwtSecret string, jwtHours int) UserService {
	return &userService{repo: r, jwtSecret: jwtSecret, jwtHours: jwtHours}
}

func (s *userService) Register(username, email, password string) (*models.User, error) {
	if username == "" || email == "" || password == "" {
		return nil, errors.New("missing fields")
	}
	// check duplicate email
	exists, err := s.repo.CheckEmailExists(email)
	if err != nil { return nil, err }
	if exists { return nil, errors.New("email already in use") }

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil { return nil, err }
	u := &models.User{Username: username, Email: email, PasswordHash: string(hash)}
	if err := s.repo.CreateUser(u); err != nil { return nil, err }
	return u, nil
}

func (s *userService) Login(email, password string) (string, *models.User, error) {
	if email == "" || password == "" { return "", nil, errors.New("missing credentials") }
	u, err := s.repo.GetUserByEmail(email)
	if err != nil { return "", nil, err }
	if u == nil { return "", nil, errors.New("invalid credentials") }
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("invalid credentials")
	}
	// generate token
	tok, err := auth.GenerateToken(u.ID, s.jwtSecret, s.jwtHours)
	if err != nil { return "", nil, err }
	return tok, u, nil
}
