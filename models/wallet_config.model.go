package models

import (
	"time"
)

type WalletConfig struct {
	ID          string          `json:"id" gorm:"column:id"`
	Endpoint    string          `json:"endpoint" gorm:"column:endpoint"`
	AccessToken string          `json:"access_token" gorm:"column:access_token"`
	CreatedAt   *time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   *time.Time      `json:"updated_at" gorm:"column:updated_at"`
}

func (m WalletConfig) TableName() string {
	return "wallet_configs"
}
