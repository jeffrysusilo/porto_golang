package routes

import (
	"crypto-wallet-service/internal/handlers"
	"crypto-wallet-service/internal/middleware"

	"github.com/gin-gonic/gin"
)


func SetupRoutes(
	router *gin.Engine,
	authHandler *handlers.AuthHandler,
	walletHandler *handlers.WalletHandler,
	transactionHandler *handlers.TransactionHandler,
) {
	
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	
	api := router.Group("/api")
	{
		
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware())
		{
		
			user := protected.Group("/user")
			{
				user.GET("/me", authHandler.GetMe)
			}

		
			wallet := protected.Group("/wallet")
			{
				wallet.GET("", walletHandler.GetWallet)
				wallet.POST("/deposit", walletHandler.Deposit)
				wallet.POST("/withdraw", walletHandler.Withdraw)
			}

			
			protected.GET("/transactions", transactionHandler.GetTransactions)
		}
	}
}
