package models

import (
	"ssi-gitlab.teda.th/ssi/core/utils"
	"gorm.io/gorm"
	"time"
)

type Organization struct {
	ID          string          `json:"id" gorm:"column:id"`
	JuristicID  string          `json:"juristic_id" gorm:"column:juristic_id"`
	Name        string          `json:"name" gorm:"column:name"`
	DIDAddress  *string         `json:"did_address" gorm:"column:did_address"`
	EncryptedID *string         `json:"encrypted_id" gorm:"column:encrypted_id"`
	CreatedAt   *time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   *time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt   *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
}

func (receiver Organization) TableName() string {
	return "organizations"
}

func NewOrganization() *Organization {
	return &Organization{
		ID:        utils.GetUUID(),
		CreatedAt: utils.GetCurrentDateTime(),
		UpdatedAt: utils.GetCurrentDateTime(),
	}
}
