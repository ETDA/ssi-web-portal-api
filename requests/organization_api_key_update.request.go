package requests

import (
	core "ssi-gitlab.teda.th/ssi/core"
)

type OrganizationAPIKeyUpdate struct {
	core.BaseValidator
	Name *string `json:"name"`
}

func (r *OrganizationAPIKeyUpdate) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.Name, "name"))

	return r.Error()
}
