package models

import (
	"gorm.io/gorm"
	"time"
)

type Key struct {
	ID                  string          `json:"id"`
	PublicKey           string          `json:"public_key"`
	PrivateKeyEncrypted string          `json:"private_key_encrypted"`
	Type                string          `json:"type"`
	CreatedAt           *time.Time      `json:"created_at"`
	UpdatedAt           *time.Time      `json:"updated_at"`
	DeletedAt           *gorm.DeletedAt `json:"deleted_at,omitempty"`
}
