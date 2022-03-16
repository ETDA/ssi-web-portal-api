package models

import (
	"encoding/json"
	"time"
)

type VCSchemaHistory struct {
	SchemaID   string           `json:"schema_id" gorm:"column:schema_id"`
	SchemaBody *json.RawMessage `json:"schema_body" gorm:"column:schema_body"`
	Schema     *VCSchema        `json:"schema,omitempty" gorm:"foreignKey:SchemaID"`
	Version    string           `json:"version" gorm:"version"`
	CreatedAt  *time.Time       `json:"created_at" gorm:"column:created_at"`
	UpdatedAt  *time.Time       `json:"updated_at" gorm:"column:updated_at"`
}

func (receiver VCSchemaHistory) TableName() string {
	return "vc_schema_histories"
}
