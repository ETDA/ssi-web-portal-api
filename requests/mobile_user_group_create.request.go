package requests

import (
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
)

type MobileUserGroupCreate struct {
	core.BaseValidator
	Name *string `json:"name"`
}

func (r *MobileUserGroupCreate) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.Name, "name"))
	r.Must(r.IsStrUnique(ctx, r.Name, models.MobileUserGroup{}.TableName(), "name", "", "name"))
	return r.Error()
}
