package main

import (
	"log"

	"github.com/shoelfikar/voucher-management-system/internal/config"
	"github.com/shoelfikar/voucher-management-system/internal/delivery/http"
	"github.com/shoelfikar/voucher-management-system/internal/delivery/http/handler"
	"github.com/shoelfikar/voucher-management-system/internal/delivery/http/middleware"
	"github.com/shoelfikar/voucher-management-system/internal/domain/entity"
	"github.com/shoelfikar/voucher-management-system/internal/repository"
	"github.com/shoelfikar/voucher-management-system/internal/service"
	"github.com/shoelfikar/voucher-management-system/pkg/database"
	"github.com/shoelfikar/voucher-management-system/pkg/jwt"
)

func main() {
	log.Println("Loading configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	log.Println("Connecting to database...")
	db, err := database.NewPostgresDatabase(&cfg.Database)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Running database migrations...")
	err = db.AutoMigrate(&entity.User{}, &entity.Voucher{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	log.Println("Initializing JWT service...")
	jwtService := jwt.NewJWTService(cfg.JWT.Secret, cfg.JWT.Expiration)

	log.Println("Initializing repositories...")
	userRepo := repository.NewUserRepository(db)
	voucherRepo := repository.NewVoucherRepository(db)

	log.Println("Initializing services...")
	authService := service.NewAuthService(userRepo, jwtService)
	voucherService := service.NewVoucherService(voucherRepo)

	log.Println("Initializing handlers...")
	authHandler := handler.NewAuthHandler(authService)
	voucherHandler := handler.NewVoucherHandler(voucherService)

	log.Println("Initializing middleware...")
	authMiddleware := middleware.AuthMiddleware(jwtService)
	corsMiddleware := middleware.CORSMiddleware(cfg.CORS.AllowedOrigins)

	log.Println("Setting up router...")
	router := http.SetupRouter(
		authHandler,
		voucherHandler,
		authMiddleware,
		corsMiddleware,
	)

	serverAddr := ":" + cfg.Server.Port
	log.Printf("Server starting on port %s (mode: %s)", cfg.Server.Port, cfg.Server.Mode)
	log.Printf("Health check: http://localhost%s/health", serverAddr)
	log.Printf("API endpoint: http://localhost%s/api/v1", serverAddr)

	if err := router.Run(serverAddr); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
