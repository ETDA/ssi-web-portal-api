package models

import (
	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"time"
)

type RequestedVP struct {
	ID          string     `json:"id" gorm:"column:id"`
	Name        string     `json:"name" gorm:"column:name"`
	Status      string     `json:"status" gorm:"column:status"`
	CreatorID   string     `json:"-" gorm:"column:creator_id"`
	Creator     *User      `json:"creator,omitempty" gorm:"foreignKey:CreatorID;References:ID"`
	Schemacount int64      `json:"schema_count" gorm:"column:schema_count"`
	QRCodeID    string     `json:"-" gorm:"column:qr_code_id"`
	CreatedAt   *time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt   *time.Time `json:"updated_at" gorm:"column:updated_at"`
}
type RequestedVPSchemaType struct {
	ID            string     `json:"id" gorm:"column:id"`
	RequestedVPID string     `json:"requested_vp_id" gorm:"requested_vp_id"`
	SchemaType    string     `json:"schema_type" gorm:"column:schema_type"`
	IsRequired    bool       `json:"is_required" gorm:"column:is_required"`
	Noted         *string    `json:"noted" gorm:"column:noted"`
	CreatedAt     *time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt     *time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (m RequestedVP) TableName() string {
	return "requested_vps"
}
func (m RequestedVPSchemaType) TableName() string {
	return "requested_vps_schema_types"
}
func NewRequestedVP() *RequestedVP {
	id := utils.GetUUID()
	createdAt := utils.GetCurrentDateTime()
	return &RequestedVP{
		ID:        id,
		Status:    consts.VPStatusActive,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
}

func NewRequestedVPSchemaType() *RequestedVPSchemaType {
	id := utils.GetUUID()
	createdAt := utils.GetCurrentDateTime()
	return &RequestedVPSchemaType{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
}
