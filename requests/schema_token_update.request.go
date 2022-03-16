package requests

import (
	"strings"

	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	core "ssi-gitlab.teda.th/ssi/core"
)

type SchemaTokenUpdate struct {
	core.BaseValidator
	Name *string `json:"name"`
	Role *string `json:"role"`
}

func (r *SchemaTokenUpdate) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.Name, "name"))
	if r.Must(r.IsRequired(r.Role, "role")) {
		r.Must(r.IsStrIn(r.Role, strings.Join([]string{consts.SchemaTokenAdminRole, consts.SchemaTokenReadWriteRole, consts.SchemaTokenReadRole}, "|"), "name"))
	}
	return r.Error()
}
