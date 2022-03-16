package requests

import core "ssi-gitlab.teda.th/ssi/core"

type VCSchemaRepositoryUpdate struct {
	core.BaseValidator
	Name     *string `json:"name"`
	Endpoint *string `json:"endpoint"`
	Token    *string `json:"token"`
}

func (r *VCSchemaRepositoryUpdate) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.Name, "name"))
	if r.Must(r.IsStrRequired(r.Endpoint, "endpoint")) {
		r.Must(r.IsURL(r.Endpoint, "endpoint"))
	}
	r.Must(r.IsStrRequired(r.Token, "token"))
	return r.Error()
}
