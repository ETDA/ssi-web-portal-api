package requests

import (
	"encoding/json"

	core "ssi-gitlab.teda.th/ssi/core"
)

type SchemaCreate struct {
	core.BaseValidator
	SchemaName *string          `json:"schema_name"`
	SchemaType *string          `json:"schema_type"`
	SchemaBody *json.RawMessage `json:"schema_body"`
}

func (r *SchemaCreate) Valid(ctx core.IContext) core.IError {

	r.Must(r.IsStrRequired(r.SchemaName, "schema_name"))
	r.Must(r.IsStrRequired(r.SchemaType, "schema_type"))
	if r.Must(r.IsRequired(r.SchemaBody, "schema_body")) {
		r.Must(r.IsJSONStrPathRequired(r.SchemaBody, "$schema", "schema_body.$schema"))
	}

	return r.Error()
}
