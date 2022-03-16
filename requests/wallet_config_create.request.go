package requests

import (
	core "ssi-gitlab.teda.th/ssi/core"
)

type WalletConfigCreate struct {
	core.BaseValidator
	Endpoint    *string `json:"endpoint"`
	AccessToken *string `json:"access_token"`
}

func (r *WalletConfigCreate) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.Endpoint, "endpoint"))
	r.Must(r.IsURL(r.Endpoint, "endpoint"))
	r.Must(r.IsStrRequired(r.AccessToken, "access_token"))

	return r.Error()
}
