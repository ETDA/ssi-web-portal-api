package models

import "time"

type MobileUser struct {
	ID        string     `json:"id"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

func (m MobileUser) TableName() string {
	return "mobile_users"
}
