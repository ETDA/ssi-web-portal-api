package requests

import core "ssi-gitlab.teda.th/ssi/core"

type VPVerify struct {
	core.BaseValidator
	JWT *string `json:"jwt"`
}

func (r *VPVerify) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.JWT, "jwt"))
	return r.Error()
}
