package models

import (
	"ssi-gitlab.teda.th/ssi/core/utils"
	"gorm.io/gorm"
	"time"
)

type OrganizationAPIKey struct {
	ID             string          `json:"id" gorm:"id"`
	OrganizationID string          `json:"organization_id" gorm:"organization_id"`
	Name           string          `json:"name" gorm:"name"`
	Key            string          `json:"key" gorm:"key"`
	Read           bool            `json:"read" gorm:"read"`
	Write          bool            `json:"write" gorm:"write"`
	CreatedAt      *time.Time      `json:"created_at" gorm:"created_at"`
	UpdatedAt      *time.Time      `json:"updated_at" gorm:"updated_at"`
	DeletedAt      *gorm.DeletedAt `json:"deleted_at" gorm:"deleted_at"`
}

func (receiver OrganizationAPIKey) TableName() string {
	return "organization_api_keys"
}

func NewAPIKey(name string) *OrganizationAPIKey {
	id := utils.GetUUID()
	createdAt := utils.GetCurrentDateTime()

	return &OrganizationAPIKey{
		ID:        id,
		Name:      name,
		Key:       utils.NewSha256(id + createdAt.String()),
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
}
