package requests

import core "ssi-gitlab.teda.th/ssi/core"

type KeyUpload struct {
	core.BaseValidator
	X509Certificate *string `json:"x509_certificate"`
	X509Key         *string `json:"x509_key"`
}

func (r *KeyUpload) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.X509Certificate, "x509_certificate"))
	r.Must(r.IsStrRequired(r.X509Key, "x509_key"))
	return r.Error()
}
