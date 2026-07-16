package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/marcell322/api-test-studio/internal/adapters/auth"
	"github.com/marcell322/api-test-studio/internal/adapters/persistence"
	"github.com/marcell322/api-test-studio/internal/config"
	"github.com/marcell322/api-test-studio/internal/server"
	"github.com/marcell322/api-test-studio/internal/domain/models"
	"github.com/marcell322/api-test-studio/internal/usecase"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(sqlite.Open(cfg.DBPath), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	// automigrate minimal models
	db.AutoMigrate(&models.User{})

	// initialize repos & services
	userRepo := persistence.NewGormUserRepository(db)
	userSvc := usecase.NewUserService(userRepo)

	// router
	r := server.NewRouter(cfg, userSvc)

	// run
	srv := &httpServer{engine: r}
	log.Printf("listening on %s", cfg.Port)
	srv.Run(cfg.Port)

	// keep main alive briefly in case of background cleanup
	time.Sleep(100 * time.Millisecond)
}

// tiny wrapper to avoid importing net/http in multiple places
type httpServer struct{ engine *gin.Engine }
func (s *httpServer) Run(addr string) { s.engine.Run(addr) }
