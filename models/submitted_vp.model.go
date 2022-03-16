package models

import (
	"time"

	"ssi-gitlab.teda.th/ssi/core/utils"
)

type SubmittedVP struct {
	ID            string       `json:"id" gorm:"column:id"`
	RequestedVPID string       `json:"requested_vp_id" gorm:"column:requested_vp_id"`
	RequestedVP   *RequestedVP `json:"requested_vp" gorm:"foreignKey:RequestedVPID;References:ID"`
	Holder        string       `json:"holder" gorm:"column:holder"`
	JWT           string       `json:"jwt" gorm:"column:jwt"`
	Tags          string       `json:"tags" gorm:"column:tags"`
	DocumentCount int64        `json:"document_count" gorm:"column:document_count"`
	Verify        bool         `json:"verify" gorm:"column:verify"`
	CreatedAt     *time.Time   `json:"created_at" gorm:"created_at"`
	UpdatedAt     *time.Time   `json:"updated_at" gorm:"updated_at"`
}

func (m SubmittedVP) TableName() string {
	return "submitted_vps"
}

func NewSubmittedVP() *SubmittedVP {
	now := utils.GetCurrentDateTime()
	return &SubmittedVP{
		ID:        utils.GetUUID(),
		CreatedAt: now,
		UpdatedAt: now,
	}
}
