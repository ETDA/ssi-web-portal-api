package requests

import (
	"fmt"

	"gitlab.finema.co/finema/etda/web-portal-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
)

type RequestVPCancelList struct {
	core.BaseValidator
	IDs []string `json:"ids"`
}

func (r *RequestVPCancelList) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsRequiredArray(r.IDs, "ids"))
	for index, id := range r.IDs {
		r.Must(r.IsStrUnique(ctx, &id, (&models.RequestedVP{}).TableName(), "id", id, fmt.Sprintf("ids[%v]", index)))
	}
	return r.Error()
}
