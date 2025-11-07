package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Wallet struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	Currency  string    `gorm:"type:varchar(10);not null" json:"currency"` // BTC, ETH, USDT, IDR
	Balance   float64   `gorm:"type:numeric(18,8);not null;default:0" json:"balance"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	User      User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
}


func (w *Wallet) BeforeCreate(tx *gorm.DB) error {
	if w.ID == uuid.Nil {
		w.ID = uuid.New()
	}
	return nil
}


type WalletWithPrice struct {
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
	PriceIDR float64 `json:"price_idr"`
	ValueIDR float64 `json:"value_idr"`
}


type PortfolioResponse struct {
	Assets        []WalletWithPrice `json:"assets"`
	TotalValueIDR float64           `json:"total_value_idr"`
}
