package handlers

import (
	"errors"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/marcell322/api-test-studio/internal/config"
	"github.com/marcell322/api-test-studio/internal/usecase"
)

type Handlers struct {
	UserSvc       usecase.UserService
	CollectionSvc usecase.CollectionService
	RequestSvc    usecase.SavedRequestService
	Cfg           *config.Config
}

func NewHandlers(us usecase.UserService, cs usecase.CollectionService, rs usecase.SavedRequestService, cfg *config.Config) *Handlers {
	return &Handlers{UserSvc: us, CollectionSvc: cs, RequestSvc: rs, Cfg: cfg}
}

// Register creates a new user account
// POST /api/register
func (h *Handlers) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request payload",
		})
		return
	}

	if err := validateRegister(req.Username, req.Email, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

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

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "invalid request payload",
		})
		return
	}

	if err := validateLogin(req.Email, req.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": err.Error(),
		})
		return
	}

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

// --- helpers shared by resource handlers ---

func parseIDParam(c *gin.Context) (uint, error) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id), nil
}

// collectionErrorStatus maps usecase-level errors to HTTP status + message.
func collectionErrorStatus(err error) (int, string) {
	switch {
	case errors.Is(err, usecase.ErrNotFound):
		return http.StatusNotFound, "collection not found"
	case errors.Is(err, usecase.ErrForbidden):
		return http.StatusForbidden, "not allowed to access this collection"
	default:
		return http.StatusInternalServerError, "internal error"
	}
}

// --- collections ---

// ListCollections
// GET /api/collections
func (h *Handlers) ListCollections(c *gin.Context) {
	userID := c.GetUint("userID")
	cols, err := h.CollectionSvc.List(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to list collections"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": cols})
}

// CreateCollection
// POST /api/collections
func (h *Handlers) CreateCollection(c *gin.Context) {
	userID := c.GetUint("userID")

	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request payload"})
		return
	}

	col, err := h.CollectionSvc.Create(userID, req.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": col})
}

// GetCollection
// GET /api/collections/:id
func (h *Handlers) GetCollection(c *gin.Context) {
	userID := c.GetUint("userID")

	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid id"})
		return
	}

	col, err := h.CollectionSvc.Get(userID, id)
	if err != nil {
		status, msg := collectionErrorStatus(err)
		c.JSON(status, gin.H{"success": false, "message": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": col})
}

// UpdateCollection (rename)
// PUT /api/collections/:id
func (h *Handlers) UpdateCollection(c *gin.Context) {
	userID := c.GetUint("userID")

	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid id"})
		return
	}

	var req struct {
		Name string `json:"name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request payload"})
		return
	}

	col, err := h.CollectionSvc.Rename(userID, id, req.Name)
	if err != nil {
		if err.Error() == "name is required" {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
			return
		}
		status, msg := collectionErrorStatus(err)
		c.JSON(status, gin.H{"success": false, "message": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": col})
}

// DeleteCollection
// DELETE /api/collections/:id
func (h *Handlers) DeleteCollection(c *gin.Context) {
	userID := c.GetUint("userID")

	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid id"})
		return
	}

	if err := h.CollectionSvc.Delete(userID, id); err != nil {
		status, msg := collectionErrorStatus(err)
		c.JSON(status, gin.H{"success": false, "message": msg})
		return
	}

	c.Status(http.StatusNoContent)
}

// --- saved requests ---

type savedRequestPayload struct {
	CollectionID uint              `json:"collection_id"`
	Name         string            `json:"name"`
	Method       string            `json:"method"`
	URL          string            `json:"url"`
	Headers      map[string]string `json:"headers"`
	Body         string            `json:"body"`
}

func requestErrorStatus(err error) (int, string) {
	switch {
	case errors.Is(err, usecase.ErrNotFound):
		return http.StatusNotFound, "not found"
	case errors.Is(err, usecase.ErrForbidden):
		return http.StatusForbidden, "not allowed to access this resource"
	default:
		return http.StatusInternalServerError, "internal error"
	}
}

// ListRequests
// GET /api/requests
// Optional query param: ?collection_id=123 to filter by collection
func (h *Handlers) ListRequests(c *gin.Context) {
	userID := c.GetUint("userID")

	if cidStr := c.Query("collection_id"); cidStr != "" {
		cid, err := strconv.ParseUint(cidStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid collection_id"})
			return
		}
		reqs, err := h.RequestSvc.ListByCollection(userID, uint(cid))
		if err != nil {
			status, msg := requestErrorStatus(err)
			c.JSON(status, gin.H{"success": false, "message": msg})
			return
		}
		c.JSON(http.StatusOK, gin.H{"success": true, "data": reqs})
		return
	}

	reqs, err := h.RequestSvc.List(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "failed to list requests"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": reqs})
}

// CreateRequest
// POST /api/requests
func (h *Handlers) CreateRequest(c *gin.Context) {
	userID := c.GetUint("userID")

	var req savedRequestPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request payload"})
		return
	}

	sr, err := h.RequestSvc.Create(userID, req.CollectionID, req.Name, req.Method, req.URL, req.Headers, req.Body)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) || errors.Is(err, usecase.ErrForbidden) {
			status, msg := requestErrorStatus(err)
			c.JSON(status, gin.H{"success": false, "message": msg})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"success": true, "data": sr})
}

// GetRequest
// GET /api/requests/:id
func (h *Handlers) GetRequest(c *gin.Context) {
	userID := c.GetUint("userID")

	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid id"})
		return
	}

	sr, err := h.RequestSvc.Get(userID, id)
	if err != nil {
		status, msg := requestErrorStatus(err)
		c.JSON(status, gin.H{"success": false, "message": msg})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": sr})
}

// UpdateRequest
// PUT /api/requests/:id
func (h *Handlers) UpdateRequest(c *gin.Context) {
	userID := c.GetUint("userID")

	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid id"})
		return
	}

	var req savedRequestPayload
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid request payload"})
		return
	}

	sr, err := h.RequestSvc.Update(userID, id, req.Name, req.Method, req.URL, req.Headers, req.Body)
	if err != nil {
		if errors.Is(err, usecase.ErrNotFound) || errors.Is(err, usecase.ErrForbidden) {
			status, msg := requestErrorStatus(err)
			c.JSON(status, gin.H{"success": false, "message": msg})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "data": sr})
}

// DeleteRequest
// DELETE /api/requests/:id
func (h *Handlers) DeleteRequest(c *gin.Context) {
	userID := c.GetUint("userID")

	id, err := parseIDParam(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid id"})
		return
	}

	if err := h.RequestSvc.Delete(userID, id); err != nil {
		status, msg := requestErrorStatus(err)
		c.JSON(status, gin.H{"success": false, "message": msg})
		return
	}

	c.Status(http.StatusNoContent)
}

// --- history endpoints (placeholder stubs, unchanged for now) ---

func (h *Handlers) ListHistory(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": []interface{}{}})
}

func (h *Handlers) GetHistoryItem(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"success": true, "data": nil})
}

func (h *Handlers) DeleteHistoryItem(c *gin.Context) {
	c.JSON(http.StatusNoContent, gin.H{"success": true})
}