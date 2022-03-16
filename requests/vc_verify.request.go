package requests

import core "ssi-gitlab.teda.th/ssi/core"

type VCVerify struct {
	core.BaseValidator
	JWT *string `json:"jwt"`
}

func (r *VCVerify) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.JWT, "jwt"))
	return r.Error()
}
