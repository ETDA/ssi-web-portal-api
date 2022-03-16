package views

import (
	"gitlab.finema.co/finema/etda/web-portal-api/models"
)

type MobileUser struct {
	ID         string  `json:"id"`
	FirstName  string  `json:"first_name"`
	LastName   string  `json:"last_name"`
	DIDAddress *string `json:"did_address,omitempty"`
}

type MobileUserWithGroup struct {
	MobileUser
	Groups []models.MobileUserGroup `json:"groups"`
}

func NewMobileUserWithGroup(user *MobileUser, groups []models.MobileUserGroup) *MobileUserWithGroup {
	return &MobileUserWithGroup{
		MobileUser: *user,
		Groups:     groups}

}
