package requests

import (
	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	core "ssi-gitlab.teda.th/ssi/core"
)

type VCUpdate struct {
	core.BaseValidator
	Status *string `json:"status"`
}

func (r *VCUpdate) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.Status, "status"))
	r.Must(r.IsStrIn(r.Status, consts.VCStatusCanceled, "status"))
	return r.Error()
}
