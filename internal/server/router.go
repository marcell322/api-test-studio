package server

import (
	"github.com/gin-gonic/gin"
	"github.com/marcell322/api-test-studio/internal/adapters/http"
	"github.com/marcell322/api-test-studio/internal/config"
	"github.com/marcell322/api-test-studio/internal/middleware"
	"github.com/marcell322/api-test-studio/internal/usecase"
)

func NewRouter(cfg *config.Config, userSvc usecase.UserService, collectionSvc usecase.CollectionService, requestSvc usecase.SavedRequestService, historySvc usecase.HistoryService) *gin.Engine {
	r := gin.Default()
	h := handlers.NewHandlers(userSvc, collectionSvc, requestSvc, historySvc, cfg)

	// public routes
	api := r.Group("/api")
	{
		// authentication
		api.POST("/register", h.Register)
		api.POST("/login", h.Login)
	}

	// protected routes (require JWT)
	protected := api.Group("")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		// user
		protected.GET("/me", h.Me)
		protected.POST("/send", h.SendRequest) // POST /api/send

		// collections
		collections := protected.Group("/collections")
		{
			collections.GET("", h.ListCollections)         // GET /api/collections
			collections.POST("", h.CreateCollection)       // POST /api/collections
			collections.GET("/:id", h.GetCollection)       // GET /api/collections/:id
			collections.PUT("/:id", h.UpdateCollection)    // PUT /api/collections/:id
			collections.DELETE("/:id", h.DeleteCollection) // DELETE /api/collections/:id
		}

		// requests
		requests := protected.Group("/requests")
		{
			requests.GET("", h.ListRequests)         // GET /api/requests
			requests.POST("", h.CreateRequest)       // POST /api/requests
			requests.GET("/:id", h.GetRequest)       // GET /api/requests/:id
			requests.PUT("/:id", h.UpdateRequest)    // PUT /api/requests/:id
			requests.DELETE("/:id", h.DeleteRequest) // DELETE /api/requests/:id
		}

		// history
		history := protected.Group("/history")
		{
			history.GET("", h.ListHistory)
			history.GET("/:id", h.GetHistoryItem)
			history.DELETE("/:id", h.DeleteHistoryItem)
		}
	}

	return r
}