package models

import (
	"time"

	"ssi-gitlab.teda.th/ssi/core/utils"
)

type NonceLock struct {
	ID        string     `json:"id" gorm:"column:id"`
	IsDone    bool       `json:"is_done" gorm:"column:is_done"`
	VCID      string     `json:"vc_id" gorm:"column:vc_id"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (m NonceLock) TableName() string {
	return "nonce_locks"
}

func NewNonceLock(vcID string) *NonceLock {

	id := utils.GetUUID()
	createdAt := utils.GetCurrentDateTime()
	return &NonceLock{
		ID:        id,
		IsDone:    false,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
		VCID:      vcID,
	}
}
