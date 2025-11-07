package handlers

import (
	"crypto-wallet-service/internal/middleware"
	"crypto-wallet-service/internal/models"
	"crypto-wallet-service/internal/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	walletService *services.WalletService
}

func NewWalletHandler(walletService *services.WalletService) *WalletHandler {
	return &WalletHandler{walletService: walletService}
}


func (h *WalletHandler) GetWallet(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	portfolio, err := h.walletService.GetPortfolio(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, portfolio)
}


func (h *WalletHandler) Deposit(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req models.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}


	validCurrencies := map[string]bool{
		"BTC": true, "ETH": true, "USDT": true, "IDR": true,
	}
	if !validCurrencies[req.Currency] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid currency. Supported: BTC, ETH, USDT, IDR"})
		return
	}

	if err := h.walletService.Deposit(userID, req.Currency, req.Amount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Deposit successful",
		"currency": req.Currency,
		"amount":   req.Amount,
	})
}


func (h *WalletHandler) Withdraw(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var req models.TransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	
	validCurrencies := map[string]bool{
		"BTC": true, "ETH": true, "USDT": true, "IDR": true,
	}
	if !validCurrencies[req.Currency] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid currency. Supported: BTC, ETH, USDT, IDR"})
		return
	}

	if err := h.walletService.Withdraw(userID, req.Currency, req.Amount); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Withdrawal successful",
		"currency": req.Currency,
		"amount":   req.Amount,
	})
}
