package views

import (
	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
)

type MobileUserGroup struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	UserCount int64  `json:"user_count"`
	IsStatic  bool   `json:"is_static"`
}

type MobileUserGroupList struct {
	Groups []MobileUserGroup `json:"groups"`
}

func NewMobileUserGroup(group *models.MobileUserGroup, userCount int64) *MobileUserGroup {
	isStatic := false
	id := group.ID
	for _, group := range consts.MobileUserGroup {
		if group == id {
			isStatic = true
			break
		}
	}
	return &MobileUserGroup{
		ID:        id,
		Name:      group.Name,
		UserCount: userCount,
		IsStatic:  isStatic,
	}
}
func NewMobileUserGroupList(groups []models.MobileUserGroup, groupUserCount []int64) []MobileUserGroup {
	groupList := make([]MobileUserGroup, 0)
	for index, group := range groups {
		groupView := NewMobileUserGroup(&group, groupUserCount[index])
		groupList = append(groupList, *groupView)
	}
	return groupList
}
