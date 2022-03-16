package requests

import (
	"strings"

	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	core "ssi-gitlab.teda.th/ssi/core"
)

type UserUpdate struct {
	core.BaseValidator
	Role *string `json:"role"`
}

func (r *UserUpdate) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrIn(r.Role, strings.Join(consts.UserRoles, "|"), "role"))
	return r.Error()
}
