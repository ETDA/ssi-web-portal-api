package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

type VCSchema struct {
	ID         string           `json:"id" gorm:"id"`
	SchemaName string           `json:"schema_name" gorm:"schema_name"`
	SchemaType string           `json:"schema_type" gorm:"schema_type"`
	SchemaBody *json.RawMessage `json:"schema_body" gorm:"schema_body"`
	Version    string           `json:"version" gorm:"version"`
	Permission string           `json:"permission"`
	CreatedBy  string           `json:"created_by" gorm:"created_by"`
	CreatedAt  *time.Time       `json:"created_at" gorm:"created_at"`
	UpdatedAt  *time.Time       `json:"updated_at" gorm:"column:updated_at"`
	DeletedAt  *gorm.DeletedAt  `json:"deleted_at,omitempty" gorm:"deleted_at"`
}

func (receiver VCSchema) TableName() string {
	return "vc_schemas"
}
