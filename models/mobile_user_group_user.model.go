package models

import (
	"ssi-gitlab.teda.th/ssi/core/utils"
	"time"
)

type MobileUserGroupUser struct {
	ID                string           `json:"id"`
	MobileUserID      string           `json:"mobile_user_id" gorm:"mobile_user_id"`
	MobileUserGroupID string           `json:"mobile_user_group_id,omitempty" gorm:"mobile_user_group_id"`
	MobileUserGroup   *MobileUserGroup `json:"group" gorm:"foreignKey:MobileUserGroupID"`
	CreatedAt         *time.Time       `json:"created_at" gorm:"column:created_at"`
	UpdatedAt         *time.Time       `json:"updated_at" gorm:"column:updated_at"`
}

func (m MobileUserGroupUser) TableName() string {
	return "mobile_user_group_users"
}

func NewMobileUserGroupUser() *MobileUserGroupUser {
	id := utils.GetUUID()
	createdAt := utils.GetCurrentDateTime()
	return &MobileUserGroupUser{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: createdAt,
	}
}
