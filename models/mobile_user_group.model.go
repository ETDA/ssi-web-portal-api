package models

import (
	"ssi-gitlab.teda.th/ssi/core/utils"
	"time"
)

type MobileUserGroup struct {
	ID        string     `json:"id" gorm:"column:id"`
	Name      string     `json:"name" gorm:"column:name"`
	CreatedAt *time.Time `json:"created_at" gorm:"column:created_at"`
	UpdatedAt *time.Time `json:"updated_at" gorm:"column:updated_at"`
}

func (m MobileUserGroup) TableName() string {
	return "mobile_user_groups"
}

func NewMobileUserGroup() *MobileUserGroup {
	id := utils.GetUUID()
	createdAt := utils.GetCurrentDateTime()
	return &MobileUserGroup{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
}
