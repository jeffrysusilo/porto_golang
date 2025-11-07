package repository

import (
	"crypto-wallet-service/internal/models"
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type WalletRepository interface {
	Create(wallet *models.Wallet) error
	FindByUserIDAndCurrency(userID uuid.UUID, currency string) (*models.Wallet, error)
	FindAllByUserID(userID uuid.UUID) ([]models.Wallet, error)
	Update(wallet *models.Wallet) error
	UpdateBalance(walletID uuid.UUID, newBalance float64) error
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db: db}
}

func (r *walletRepository) Create(wallet *models.Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *walletRepository) FindByUserIDAndCurrency(userID uuid.UUID, currency string) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.Where("user_id = ? AND currency = ?", userID, currency).First(&wallet).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Return nil if not found (not an error for wallet creation)
		}
		return nil, err
	}
	return &wallet, nil
}

func (r *walletRepository) FindAllByUserID(userID uuid.UUID) ([]models.Wallet, error) {
	var wallets []models.Wallet
	err := r.db.Where("user_id = ?", userID).Find(&wallets).Error
	if err != nil {
		return nil, err
	}
	return wallets, nil
}

func (r *walletRepository) Update(wallet *models.Wallet) error {
	return r.db.Save(wallet).Error
}

func (r *walletRepository) UpdateBalance(walletID uuid.UUID, newBalance float64) error {
	return r.db.Model(&models.Wallet{}).
		Where("id = ?", walletID).
		Update("balance", newBalance).Error
}
