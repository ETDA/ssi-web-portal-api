package models

import (
	"time"

	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/helpers"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type User struct {
	ID             string        `json:"id" gorm:"column:id"`
	Email          string        `json:"email" gorm:"column:email"`
	Password       *string       `json:"-" gorm:"column:password"`
	FirstName      string        `json:"first_name" gorm:"column:first_name"`
	LastName       string        `json:"last_name" gorm:"column:last_name"`
	Status         string        `json:"status" gorm:"column:status"`
	Role           string        `json:"role" gorm:"column:role"`
	DateOfBirth    string        `json:"date_of_birth" gorm:"column:date_of_birth"`
	VerifyToken    *string       `json:"-" gorm:"column:verify_token"`
	OrganizationID string        `json:"organization_id" gorm:"organization_id"`
	Organization   *Organization `json:"organization,omitempty" gorm:"foreignKey:OrganizationID"`
	CreatedAt      *time.Time    `json:"created_at" gorm:"column:created_at"`
	UpdatedAt      *time.Time    `json:"updated_at" gorm:"column:updated_at"`
}

func (m User) TableName() string {
	return "users"
}

func NewUser() *User {
	id := utils.GetUUID()
	createdAt := utils.GetCurrentDateTime()
	token := utils.NewSha256(id + createdAt.String() + consts.UserStatusInActive + utils.GetUUID())
	return &User{
		ID:          id,
		Role:        string(consts.UserRoleMember),
		Status:      consts.UserStatusInActive,
		VerifyToken: &token,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}
}

func (m *User) SetPassword(password string) {
	hashPassword, _ := helpers.HashPassword(password)
	m.Password = hashPassword
}
