package handlers

import (
	"errors"
	"net/http"
	"net/mail"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/marcell322/api-test-studio/internal/config"
	"github.com/marcell322/api-test-studio/internal/usecase"
)

type Handlers struct {
	UserSvc usecase.UserService
	Cfg     *config.Config
}

func NewHandlers(us usecase.UserService, cfg *config.Config) *Handlers {
	return &Handlers{UserSvc: us, Cfg: cfg}
}

// Register creates a new user account
// POST /api/register
func (h *Handlers) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// bind and parse JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request payload",
		})
		return
	}

	// validate input
	if err := validateRegister(req.Username, req.Email, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// call service layer
	user, err := h.UserSvc.Register(req.Username, req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    user,
	})
}

// Login authenticates user and returns JWT token
// POST /api/login
func (h *Handlers) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// bind and parse JSON
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request payload",
		})
		return
	}

	// validate input
	if err := validateLogin(req.Email, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

	// call service layer
	token, user, err := h.UserSvc.Login(req.Email, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "invalid credentials",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"token": token,
			"user":  user,
		},
	})
}

// Me returns the authenticated user info
// GET /api/me (requires JWT)
func (h *Handlers) Me(c *gin.Context) {
	userIDI, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"message": "unauthenticated",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"user_id": userIDI,
		},
	})
}

// validateRegister validates registration input
func validateRegister(username, email, password string) error {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)

	if username == "" {
		return errors.New("username is required")
	}
	if len(username) < 3 {
		return errors.New("username must be at least 3 characters")
	}
	if email == "" {
		return errors.New("email is required")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("email format is invalid")
	}
	if password == "" {
		return errors.New("password is required")
	}
	if len(password) < 6 {
		return errors.New("password must be at least 6 characters")
	}
	return nil
}

// validateLogin validates login input
func validateLogin(email, password string) error {
	email = strings.TrimSpace(email)

	if email == "" {
		return errors.New("email is required")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return errors.New("email format is invalid")
	}
	if password == "" {
		return errors.New("password is required")
	}
	return nil
}

// collections endpoints (placeholder stubs)

func (h *Handlers) ListCollections(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": []interface{}{}})
}

func (h *Handlers) CreateCollection(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": nil})
}

func (h *Handlers) GetCollection(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": nil})
}

func (h *Handlers) UpdateCollection(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": nil})
}

func (h *Handlers) DeleteCollection(c *gin.Context) {
	c.JSON(http.StatusNoContent, gin.H{"success": true})
}

// requests endpoints (placeholder stubs)

func (h *Handlers) ListRequests(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": []interface{}{}})
}

func (h *Handlers) CreateRequest(c *gin.Context) {
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": nil})
}

func (h *Handlers) GetRequest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": nil})
}

func (h *Handlers) UpdateRequest(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": nil})
}

func (h *Handlers) DeleteRequest(c *gin.Context) {
	c.JSON(http.StatusNoContent, gin.H{"success": true})
}

// history endpoints (placeholder stubs)

func (h *Handlers) ListHistory(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": []interface{}{}})
}

func (h *Handlers) GetHistoryItem(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": nil})
}

func (h *Handlers) DeleteHistoryItem(c *gin.Context) {
	c.JSON(http.StatusNoContent, gin.H{"success": true})
}
