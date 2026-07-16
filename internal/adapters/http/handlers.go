package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/marcell322/api-test-studio/internal/adapters/auth"
	"github.com/marcell322/api-test-studio/internal/config"
	"github.com/marcell322/api-test-studio/internal/usecase"
)

type Handlers struct{
	UserSvc usecase.UserService
	Cfg     *config.Config
}

func NewHandlers(us usecase.UserService, cfg *config.Config) *Handlers { return &Handlers{UserSvc: us, Cfg: cfg} }

func (h *Handlers) Register(c *gin.Context) {
	var req struct{ Username, Email, Password string }
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid payload"}); return }
	u, err := h.UserSvc.Register(req.Username, req.Email, req.Password)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": err.Error()}); return }
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": u})
}

func (h *Handlers) Login(c *gin.Context) {
	var req struct{ Email, Password string }
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"success": false, "message": "invalid payload"}); return }
	u, err := h.UserSvc.Authenticate(req.Email, req.Password)
	if err != nil { c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "invalid credentials"}); return }
	token, err := auth.GenerateToken(u.ID, h.Cfg.JWTSecret, h.Cfg.JWTExpireH)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"success": false, "message": "could not generate token"}); return }
	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"token": token}})
}

func (h *Handlers) Me(c *gin.Context) {
	idI, ok := c.Get("userID")
	if !ok { c.JSON(http.StatusUnauthorized, gin.H{"success": false, "message": "unauthenticated"}); return }
	c.JSON(http.StatusOK, gin.H{"success": true, "data": gin.H{"user_id": idI}})
}
