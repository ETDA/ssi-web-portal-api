package requests

import (
	"encoding/json"

	core "ssi-gitlab.teda.th/ssi/core"
)

type SchemaUpdate struct {
	core.BaseValidator
	SchemaName *string          `json:"schema_name"`
	SchemaBody *json.RawMessage `json:"schema_body"`
	Version    *string          `json:"version"`
}

func (r *SchemaUpdate) Valid(ctx core.IContext) core.IError {

	r.Must(r.IsStrRequired(r.SchemaName, "schema_name"))
	if r.Must(r.IsRequired(r.SchemaBody, "schema_body")) {
		r.Must(r.IsJSONStrPathRequired(r.SchemaBody, "$schema", "schema_body.$schema"))

	}

	r.Must(r.IsStrRequired(r.Version, "version"))
	return r.Error()
}
