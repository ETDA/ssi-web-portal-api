package requests

import core "ssi-gitlab.teda.th/ssi/core"

type UserPasswordChange struct {
	core.BaseValidator
	CurrentPassword *string `json:"current_password"`
	NewPassword     *string `json:"new_password"`
}

func (r *UserPasswordChange) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsRequired(r.CurrentPassword, "current_password"))
	r.Must(r.IsRequired(r.NewPassword, "new_password"))
	r.Must(r.IsStrMin(r.NewPassword, 8, "new_password"))
	return r.Error()
}
