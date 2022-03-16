package requests

import (
	"fmt"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
)

type RequestedVPCreate struct {
	core.BaseValidator
	Name       *string               `json:"name"`
	SchemaList []RequestedSchemaList `json:"schema_list"`
}
type RequestedSchemaList struct {
	SchemaType *string `json:"schema_type"`
	IsRequired *bool   `json:"is_required"`
	Noted      *string `json:"noted"`
}

func (r *RequestedVPCreate) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.Name, "name"))
	r.Must(r.IsStrUnique(ctx, r.Name, models.RequestedVP{}.TableName(), "name", "", "name"))
	if r.Must(r.IsRequiredArray(r.SchemaList, "schema_list")) {
		for index, requestedSchema := range r.SchemaList {
			r.Must(r.IsStrRequired(requestedSchema.SchemaType, fmt.Sprintf("schema_list[%v].schema_type", index)))
			r.Must(r.IsRequired(requestedSchema.IsRequired, fmt.Sprintf("schema_list[%v].is_required", index)))
		}
	}
	return r.Error()
}
