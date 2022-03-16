package requests

import core "ssi-gitlab.teda.th/ssi/core"

type VCSchemaRepositoryCreate struct {
	core.BaseValidator
	SchemaConfigs []VCSchemaRepositoryConfig `json:"schema_configs"`
}
type VCSchemaRepositoryConfig struct {
	Name     *string `json:"name"`
	Endpoint *string `json:"endpoint"`
	Token    *string `json:"token"`
}

func (r *VCSchemaRepositoryCreate) Valid(ctx core.IContext) core.IError {
	if r.Must(r.IsRequiredArray(r.SchemaConfigs, "schema_configs")) {

		for _, schemaConfig := range r.SchemaConfigs {
			r.Must(r.IsStrRequired(schemaConfig.Name, "name"))
			if r.Must(r.IsStrRequired(schemaConfig.Endpoint, "endpoint")) {
				r.Must(r.IsURL(schemaConfig.Endpoint, "endpoint"))
			}
			r.Must(r.IsStrRequired(schemaConfig.Token, "token"))
		}
	}
	return r.Error()
}
