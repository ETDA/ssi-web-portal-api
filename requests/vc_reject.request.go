package requests

import core "ssi-gitlab.teda.th/ssi/core"

type VCReject struct {
	core.BaseValidator
	Reason *string `json:"reason"`
}

func (r *VCReject) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.Reason, "reason"))
	return r.Error()
}
