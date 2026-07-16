package server

import (
	"github.com/gin-gonic/gin"
	"github.com/marcell322/api-test-studio/internal/adapters/http"
	"github.com/marcell322/api-test-studio/internal/config"
	"github.com/marcell322/api-test-studio/internal/middleware"
	"github.com/marcell322/api-test-studio/internal/usecase"
)

func NewRouter(cfg *config.Config, userSvc usecase.UserService) *gin.Engine {
	r := gin.Default()
	h := handlers.NewHandlers(userSvc, cfg)

	api := r.Group("/api")
	api.POST("/register", h.Register)
	api.POST("/login", h.Login)

	private := api.Group("")
	private.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	private.GET("/me", h.Me)

	return r
}
