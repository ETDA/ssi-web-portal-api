package requests

import core "ssi-gitlab.teda.th/ssi/core"

type UserForgotPassword struct {
	core.BaseValidator
	Email *string `json:"email"`
}

func (r *UserForgotPassword) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.Email, "email"))
	r.Must(r.IsEmail(r.Email, "email"))
	return r.Error()
}
