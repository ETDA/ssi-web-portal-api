package requests

import core "ssi-gitlab.teda.th/ssi/core"

type VCSignRequest struct {
	core.BaseValidator
	SchemaName        *string                       `json:"schema_name"`
	Signer            *string                       `json:"signer"` //Not required, if provided one notify moblie client to sign. If not provided the web portal server sign from key in key service.
	Holder            *string                       `json:"holder"`
	CredentialSubject core.Map                      `json:"credentialSubject"`
	CredentialSchema  *VCJWTMessageCredentialSchema `json:"credentialSchema"`
}

func (r *VCSignRequest) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.SchemaName, "schema_name"))
	r.Must(r.IsStrRequired(r.Holder, "holder"))
	r.Must(r.IsRequired(r.CredentialSubject, "credentialSubject"))
	if r.Must(r.IsRequired(r.CredentialSchema, "credentialSchema")) {
		r.Must(r.IsStrRequired(r.CredentialSchema.ID, "credentialSchema.id"))
		r.Must(r.IsStrRequired(r.CredentialSchema.Type, "credentialSchema.type"))
	}
	return r.Error()
}
