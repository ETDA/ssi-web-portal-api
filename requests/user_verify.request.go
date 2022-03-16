package requests

import core "ssi-gitlab.teda.th/ssi/core"

type UserVerify struct {
	core.BaseValidator
	Token    *string `json:"token"`
	Password *string `json:"password"`
}

func (r *UserVerify) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.Token, "token"))
	r.Must(r.IsStrRequired(r.Password, "password"))
	r.Must(r.IsStrMin(r.Password, 8, "password"))
	return r.Error()
}
