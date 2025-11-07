package main

import (
	"crypto-wallet-service/config"
	"crypto-wallet-service/internal/handlers"
	"crypto-wallet-service/internal/repository"
	"crypto-wallet-service/internal/routes"
	"crypto-wallet-service/internal/services"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	
	cfg := config.LoadConfig()


	db, err := config.InitDatabase(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}


	redisClient, err := config.InitRedis(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to initialize Redis: %v", err)
	}


	userRepo := repository.NewUserRepository(db)
	walletRepo := repository.NewWalletRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)

	
	coinGeckoService := services.NewCoinGeckoService(redisClient)
	walletService := services.NewWalletService(walletRepo, transactionRepo, coinGeckoService, db)


	authHandler := handlers.NewAuthHandler(userRepo)
	walletHandler := handlers.NewWalletHandler(walletService)
	transactionHandler := handlers.NewTransactionHandler(transactionRepo)

	
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()


	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})


	routes.SetupRoutes(router, authHandler, walletHandler, transactionHandler)


	log.Printf("Server starting on port %s...", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
