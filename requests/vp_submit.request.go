package requests

import core "ssi-gitlab.teda.th/ssi/core"

type VPSubmit struct {
	core.BaseValidator
}

func (r *VPSubmit) Valid(ctx core.IContext) core.IError {
	return r.Error()
}
