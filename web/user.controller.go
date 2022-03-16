package web

import (
	"net/http"

	"gitlab.finema.co/finema/etda/web-portal-api/requests"
	"gitlab.finema.co/finema/etda/web-portal-api/services"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type UserController struct{}

func (a UserController) Update(c core.IHTTPContext) error {
	input := &requests.UserUpdate{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	userSvc := services.NewUserService(c)
	item, ierr := userSvc.Update(c.Param("user_id"), &services.UserUpdatePayload{
		Role: utils.GetString(input.Role),
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, item)
}

func (a UserController) Verify(c core.IHTTPContext) error {
	input := &requests.UserVerify{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	service := services.NewUserService(c)
	item, ierr := service.Verify(utils.GetString(input.Token), utils.GetString(input.Password))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusCreated, item)
}

func (a UserController) CheckVerifyToken(c core.IHTTPContext) error {

	service := services.NewUserService(c)
	ierr := service.CheckVerifyToken(c.QueryParam("token"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.NoContent(http.StatusNoContent)
}

func (a UserController) Register(c core.IHTTPContext) error {
	input := &requests.UserRegister{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	userSvc := services.NewUserService(c)
	organizationSvc := services.NewOrganizationService(c)
	organization, ierr := organizationSvc.First()
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	userRegisterUserPayload := &services.UserRegisterPayload{}
	_ = utils.Copy(userRegisterUserPayload, input)
	userRegisterUserPayload.OrganizationID = organization.ID
	item, ierr := userSvc.Register(userRegisterUserPayload)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusCreated, item)
}

func (a UserController) Pagination(c core.IHTTPContext) error {
	userSvc := services.NewUserService(c)
	items, pageResponse, ierr := userSvc.Pagination(c.GetPageOptions())
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, core.NewPagination(items, pageResponse))
}

func (a UserController) Delete(c core.IHTTPContext) error {
	userSvc := services.NewUserService(c)

	ierr := userSvc.Delete(c.Param("user_id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.NoContent(http.StatusNoContent)
}

func (a UserController) Find(c core.IHTTPContext) error {
	userSvc := services.NewUserService(c)
	item, ierr := userSvc.Find(c.Param("user_id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, item)
}

func (a UserController) ResetPassword(c core.IHTTPContext) error {

	userSvc := services.NewUserService(c)
	ierr := userSvc.ResetPassword(c.Param("user_id"))
	if ierr != nil {
		if ierr != nil {
			return c.JSON(ierr.GetStatus(), ierr.JSON())
		}

	}
	return c.NoContent(http.StatusNoContent)
}
