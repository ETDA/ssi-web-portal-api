package requests

import (
	"fmt"

	"gitlab.finema.co/finema/etda/web-portal-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
)

type MobileUserUpdate struct {
	core.BaseValidator
	GroupIDs []string `json:"group_ids"`
}

func (r *MobileUserUpdate) Valid(ctx core.IContext) core.IError {
	if r.Must(r.IsRequiredArray(r.GroupIDs, "group_ids")) {
		for index, groupID := range r.GroupIDs {
			r.Must(r.IsStrUnique(ctx, &groupID, models.MobileUserGroup{}.TableName(), "id", groupID, fmt.Sprintf("user_ids[%d]", index)))
		}
	}
	return r.Error()
}
