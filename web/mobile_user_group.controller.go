package web

import (
	"gitlab.finema.co/finema/etda/web-portal-api/requests"
	"gitlab.finema.co/finema/etda/web-portal-api/services"
	"net/http"
	"ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type MobileGroupController struct{}

func (n *MobileGroupController) Get(c core.IHTTPContext) error {
	mobileUserGroupSvc := services.NewMobileUserGroupService(c)
	mobileUserGroups, ierr := mobileUserGroupSvc.Get()
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, core.NewPagination(mobileUserGroups, &core.PageResponse{
		Total:   int64(len(mobileUserGroups)),
		Limit:   1000,
		Count:   int64(len(mobileUserGroups)),
		Page:    1,
		Q:       "",
		OrderBy: nil,
	}))
}

func (n *MobileGroupController) GroupCreate(c core.IHTTPContext) error {
	input := &requests.MobileUserGroupCreate{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}
	mobileUserGroupSvc := services.NewMobileUserGroupService(c)
	group, ierr := mobileUserGroupSvc.Create(&services.MobileUserGroupCreatePayload{
		Name: utils.GetString(input.Name),
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusCreated, group)
}
func (n *MobileGroupController) GroupUpdate(c core.IHTTPContext) error {
	input := &requests.MobileUserGroupUpdate{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}
	mobileUserGroupSvc := services.NewMobileUserGroupService(c)
	group, ierr := mobileUserGroupSvc.Update(c.Param("id"), &services.MobileUserGroupUpdatePayload{
		Name: utils.GetString(input.Name),
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusCreated, group)
}

func (n *MobileGroupController) GroupDelete(c core.IHTTPContext) error {
	mobileUserGroupSvc := services.NewMobileUserGroupService(c)
	ierr := mobileUserGroupSvc.Delete(c.Param("id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.NoContent(http.StatusNoContent)
}
func (n *MobileGroupController) AddGroupUser(c core.IHTTPContext) error {
	input := &requests.MobileUserGroupAddUser{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}
	mobileUserGroupSvc := services.NewMobileUserGroupService(c)
	_, ierr := mobileUserGroupSvc.AddUser(input.GroupIDs, input.UserIDs)

	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.NoContent(http.StatusNoContent)
}

func (n *MobileGroupController) RemoveGroupUser(c core.IHTTPContext) error {
	input := &requests.MobileUserGroupRemoveUser{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}
	mobileUserGroupSvc := services.NewMobileUserGroupService(c)
	_, ierr := mobileUserGroupSvc.RemoveUser(c.Param("id"), input.UserIDs)

	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.NoContent(http.StatusNoContent)
}
func (n *MobileGroupController) Pagination(c core.IHTTPContext) error {
	mobileUserGroupSvc := services.NewMobileUserGroupService(c)
	users, pageResponse, ierr := mobileUserGroupSvc.UserPagination(c.Param("id"), c.GetPageOptions())
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, core.NewPagination(users, pageResponse))
}
