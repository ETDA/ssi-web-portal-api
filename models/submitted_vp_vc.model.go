package models

import (
	"time"

	"ssi-gitlab.teda.th/ssi/core/utils"
)

type SubmittedVPVC struct {
	ID            string     `json:"id" gorm:"column:id"`
	SubmittedVPID string     `json:"submitted_vp_id" gorm:"column:submitted_vp_id"`
	CID           string     `json:"cid" gorm:"column:cid"`
	SchemaName    string     `json:"schema_name" gorm:"column:schema_name"`
	SchemaType    string     `json:"schema_type" gorm:"column:schema_type"`
	IssuanceDate  *time.Time `json:"issuance_date" gorm:"column:issuance_date"`
	Issuer        string     `json:"issuer" gorm:"column:issuer"`
	Holder        string     `json:"holder" gorm:"column:holder"`
	JWT           string     `json:"jwt" gorm:"column:jwt"`
	Status        string     `json:"status" gorm:"column:status"`
	Verify        bool       `json:"verify" gorm:"column:verify"`
	CreatedAt     *time.Time `json:"created_at" gorm:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at" gorm:"updated_at"`
}

func (m SubmittedVPVC) TableName() string {
	return "submitted_vp_vcs"
}
func NewSubmittedVPVC() *SubmittedVPVC {
	now := utils.GetCurrentDateTime()
	return &SubmittedVPVC{
		ID:        utils.GetUUID(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}
