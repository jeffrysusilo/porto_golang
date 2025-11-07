package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TransactionType string

const (
	TransactionTypeDeposit  TransactionType = "deposit"
	TransactionTypeWithdraw TransactionType = "withdraw"
	TransactionTypeTransfer TransactionType = "transfer"
)

type Transaction struct {
	ID        uuid.UUID       `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID       `gorm:"type:uuid;not null;index" json:"user_id"`
	Type      TransactionType `gorm:"type:varchar(20);not null" json:"type"`
	Currency  string          `gorm:"type:varchar(10);not null" json:"currency"`
	Amount    float64         `gorm:"type:numeric(18,8);not null" json:"amount"`
	PriceAt   float64         `gorm:"type:numeric(18,2)" json:"price_at"` // Harga crypto saat transaksi (dalam IDR)
	CreatedAt time.Time       `gorm:"autoCreateTime" json:"created_at"`
	User      User            `gorm:"foreignKey:UserID" json:"user,omitempty"`
}


func (t *Transaction) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}


type TransactionRequest struct {
	Currency string  `json:"currency" binding:"required"`
	Amount   float64 `json:"amount" binding:"required,gt=0"`
}
