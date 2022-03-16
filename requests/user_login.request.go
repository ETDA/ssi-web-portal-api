package requests

import core "ssi-gitlab.teda.th/ssi/core"

type UserLogin struct {
	core.BaseValidator
	Email    *string `json:"email"`
	Password *string `json:"password"`
}

func (r *UserLogin) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.Email, "email"))
	r.Must(r.IsStrRequired(r.Email, "password"))
	return r.Error()
}
