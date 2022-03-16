package requests

import (
	core "ssi-gitlab.teda.th/ssi/core"
)

type OrganizationAPIKeyCreate struct {
	core.BaseValidator
	Name  *string `json:"name"`
	Read  *bool   `json:"read"`
	Write *bool   `json:"write"`
}

func (r *OrganizationAPIKeyCreate) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.Name, "name"))
	r.Must(r.Read != nil, core.RequiredM("read"))
	r.Must(r.Write != nil, core.RequiredM("write"))

	return r.Error()
}
