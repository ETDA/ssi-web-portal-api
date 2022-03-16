package requests

import (
	"strings"

	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
)

type UserRegister struct {
	core.BaseValidator
	Email       *string `json:"email"`
	FirstName   *string `json:"first_name"`
	LastName    *string `json:"last_name"`
	Role        *string `json:"role"`
	DateOfBirth *string `json:"date_of_birth"`
}

func (r *UserRegister) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.Email, "email"))
	r.Must(r.IsEmail(r.Email, "email"))
	r.Must(r.IsStrUnique(ctx, r.Email, models.User{}.TableName(), "email", "", "email"))
	r.Must(r.IsStrRequired(r.FirstName, "first_name"))
	r.Must(r.IsStrRequired(r.LastName, "last_name"))
	r.Must(r.IsStrRequired(r.Role, "role"))
	r.Must(r.IsStrIn(r.Role, strings.Join(consts.UserRoles, "|"), "role"))
	r.Must(r.IsStrRequired(r.DateOfBirth, "date_of_birth"))
	return r.Error()
}
