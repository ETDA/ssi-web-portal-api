package web

import "C"
import (
	"net/http"

	"gitlab.finema.co/finema/etda/web-portal-api/requests"
	"gitlab.finema.co/finema/etda/web-portal-api/services"
	core "ssi-gitlab.teda.th/ssi/core"
)

type MobileUserController struct{}

func (n *MobileUserController) Pagination(c core.IHTTPContext) error {
	mobileUserService := services.NewMobileUserService(c)
	users, pageResponse, ierr := mobileUserService.Pagination(c.GetPageOptions())
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, core.NewPagination(users, pageResponse))
}

func (n *MobileUserController) Find(c core.IHTTPContext) error {
	mobileUserService := services.NewMobileUserService(c)
	users, ierr := mobileUserService.FindByID(c.Param("id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, users)

}
func (n *MobileUserController) Update(c core.IHTTPContext) error {
	input := &requests.MobileUserUpdate{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}
	mobileUserGroupSvc := services.NewMobileUserGroupService(c)
	err := mobileUserGroupSvc.UpdateUser(input.GroupIDs, c.Param("id"))
	if err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}
	return c.NoContent(http.StatusNoContent)
}
