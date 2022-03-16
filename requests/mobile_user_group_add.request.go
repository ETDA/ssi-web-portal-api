package requests

import (
	"fmt"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
)

type MobileUserGroupAddUser struct {
	core.BaseValidator
	GroupIDs []string `json:"group_ids"`
	UserIDs  []string `json:"user_ids"`
}

func (r *MobileUserGroupAddUser) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsRequiredArray(r.GroupIDs, "group_ids"))
	for index, groupID := range r.GroupIDs {
		r.Must(r.IsStrUnique(ctx, &groupID, models.MobileUserGroup{}.TableName(), "id", groupID, fmt.Sprintf("user_ids[%d]", index)))
	}
	r.Must(r.IsRequiredArray(r.UserIDs, "user_ids"))
	for index, userID := range r.UserIDs {
		r.Must(r.IsStrUnique(ctx, &userID, models.MobileUser{}.TableName(), "id", userID, fmt.Sprintf("user_ids[%d]", index)))
	}
	return r.Error()
}
