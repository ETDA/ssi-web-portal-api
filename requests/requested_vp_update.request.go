package requests

import (
	"fmt"
	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	core "ssi-gitlab.teda.th/ssi/core"
)

type RequestVPUpdate struct {
	core.BaseValidator
	Status *string `json:"status"`
}

func (r *RequestVPUpdate) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.Status, "name"))
	r.Must(r.IsStrIn(r.Status, fmt.Sprintf("%s|%s|%s", consts.VPStatusActive, consts.VPStatusInActive, consts.VPStatusCancel), "status"))
	return r.Error()
}
