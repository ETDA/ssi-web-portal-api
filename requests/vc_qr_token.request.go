package requests

import core "ssi-gitlab.teda.th/ssi/core"

type VCQRTokenRequest struct {
	core.BaseValidator
	CIDs       []string `json:"cids"`
	DIDAddress *string  `json:"did_address"`
}

func (r *VCQRTokenRequest) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsRequiredArray(r.CIDs, "cids"))
	r.Must(r.IsRequired(r.DIDAddress, "did_address"))
	return r.Error()
}
