package models

import (
	"encoding/json"
	"time"

	"ssi-gitlab.teda.th/ssi/core/utils"
)

type UserAccessToken struct {
	ID        string           `json:"id" gorm:"id"`
	UserID    string           `json:"user_id" gorm:"user_id"`
	Token     string           `json:"token" gorm:"token"`
	Info      *json.RawMessage `json:"info" gorm:"info"`
	CreatedAt *time.Time       `json:"created_at" gorm:"created_at"`
}

func (receiver UserAccessToken) TableName() string {
	return "user_access_tokens"
}

func NewUserAccessToken(userID string) *UserAccessToken {
	id := utils.GetUUID()
	createdAt := utils.GetCurrentDateTime()
	info := json.RawMessage(`{}`)
	return &UserAccessToken{
		ID:        id,
		UserID:    userID,
		Token:     utils.NewSha256(id + createdAt.String()),
		CreatedAt: createdAt,
		Info:      &info,
	}
}
