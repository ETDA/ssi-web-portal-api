package models

import (
	"ssi-gitlab.teda.th/ssi/core/utils"
	"time"

	"gorm.io/gorm"
)

type VCSchemaConfig struct {
	ID          string          `json:"id" gorm:"column:id"`
	Name        string          `json:"name" gorm:"column:name"`
	Endpoint    string          `json:"endpoint" gorm:"column:endpoint"`
	AccessToken string          `json:"access_token" gorm:"column:access_token"`
	Permission  string          `json:"permission" gorm:"column:permission"`
	CreatedAt   *time.Time      `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   *time.Time      `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt   *gorm.DeletedAt `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
}

func (m VCSchemaConfig) TableName() string {
	return "vc_schema_configs"
}

func NewVCSchemaConfig(name string, endpoint string, accessToken string) *VCSchemaConfig {
	now := utils.GetCurrentDateTime()
	return &VCSchemaConfig{
		ID:          utils.GetUUID(),
		Name:        name,
		Endpoint:    endpoint,
		AccessToken: accessToken,
		Permission:  "",
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}
