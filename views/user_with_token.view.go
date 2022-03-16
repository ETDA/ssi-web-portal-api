package views

import "gitlab.finema.co/finema/etda/web-portal-api/models"

type UserWithToken struct {
	models.User
	Token string `json:"token"`
}
