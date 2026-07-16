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

	db, err := persistence.NewGormDB(cfg.DBPath)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}

	// automigrate minimal models
	if err := persistence.AutoMigrate(db, &models.User{}); err != nil {
		log.Fatalf("automigrate failed: %v", err)
	}

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
