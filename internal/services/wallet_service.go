package services

import (
	"crypto-wallet-service/internal/models"
	"crypto-wallet-service/internal/repository"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletService struct {
	walletRepo      repository.WalletRepository
	transactionRepo repository.TransactionRepository
	coinGeckoSvc    *CoinGeckoService
	db              *gorm.DB
}

func NewWalletService(
	walletRepo repository.WalletRepository,
	transactionRepo repository.TransactionRepository,
	coinGeckoSvc *CoinGeckoService,
	db *gorm.DB,
) *WalletService {
	return &WalletService{
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		coinGeckoSvc:    coinGeckoSvc,
		db:              db,
	}
}


func (s *WalletService) GetOrCreateWallet(userID uuid.UUID, currency string) (*models.Wallet, error) {
	wallet, err := s.walletRepo.FindByUserIDAndCurrency(userID, currency)
	if err != nil {
		return nil, err
	}

	if wallet == nil {
		
		wallet = &models.Wallet{
			UserID:   userID,
			Currency: currency,
			Balance:  0,
		}
		if err := s.walletRepo.Create(wallet); err != nil {
			return nil, err
		}
	}

	return wallet, nil
}

// Deposit adds funds to wallet
func (s *WalletService) Deposit(userID uuid.UUID, currency string, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get or create wallet
	wallet, err := s.GetOrCreateWallet(userID, currency)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Update balance
	newBalance := wallet.Balance + amount
	if err := s.walletRepo.UpdateBalance(wallet.ID, newBalance); err != nil {
		tx.Rollback()
		return err
	}

	// Get current price
	price, err := s.coinGeckoSvc.GetPrice(currency)
	if err != nil {
		price = 0 // Set to 0 if price fetch fails
	}

	// Create transaction record
	transaction := &models.Transaction{
		UserID:   userID,
		Type:     models.TransactionTypeDeposit,
		Currency: currency,
		Amount:   amount,
		PriceAt:  price,
	}

	if err := s.transactionRepo.Create(transaction); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// Withdraw removes funds from wallet
func (s *WalletService) Withdraw(userID uuid.UUID, currency string, amount float64) error {
	if amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	// Start transaction
	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get wallet
	wallet, err := s.walletRepo.FindByUserIDAndCurrency(userID, currency)
	if err != nil {
		tx.Rollback()
		return err
	}

	if wallet == nil {
		tx.Rollback()
		return errors.New("wallet not found")
	}

	// Check if sufficient balance
	if wallet.Balance < amount {
		tx.Rollback()
		return errors.New("insufficient balance")
	}

	// Update balance
	newBalance := wallet.Balance - amount
	if err := s.walletRepo.UpdateBalance(wallet.ID, newBalance); err != nil {
		tx.Rollback()
		return err
	}

	// Get current price
	price, err := s.coinGeckoSvc.GetPrice(currency)
	if err != nil {
		price = 0
	}

	// Create transaction record
	transaction := &models.Transaction{
		UserID:   userID,
		Type:     models.TransactionTypeWithdraw,
		Currency: currency,
		Amount:   amount,
		PriceAt:  price,
	}

	if err := s.transactionRepo.Create(transaction); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// GetPortfolio returns user's portfolio with current prices
func (s *WalletService) GetPortfolio(userID uuid.UUID) (*models.PortfolioResponse, error) {
	// Get all wallets
	wallets, err := s.walletRepo.FindAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	var assets []models.WalletWithPrice
	var totalValueIDR float64

	for _, wallet := range wallets {
		if wallet.Balance == 0 {
			continue // Skip empty wallets
		}

		// Get current price
		price, err := s.coinGeckoSvc.GetPrice(wallet.Currency)
		if err != nil {
			return nil, fmt.Errorf("failed to get price for %s: %w", wallet.Currency, err)
		}

		valueIDR := wallet.Balance * price

		assets = append(assets, models.WalletWithPrice{
			Currency: wallet.Currency,
			Balance:  wallet.Balance,
			PriceIDR: price,
			ValueIDR: valueIDR,
		})

		totalValueIDR += valueIDR
	}

	return &models.PortfolioResponse{
		Assets:        assets,
		TotalValueIDR: totalValueIDR,
	}, nil
}

// GetWallets returns all user wallets
func (s *WalletService) GetWallets(userID uuid.UUID) ([]models.Wallet, error) {
	return s.walletRepo.FindAllByUserID(userID)
}
