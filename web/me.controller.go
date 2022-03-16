package web

import (
	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	"gitlab.finema.co/finema/etda/web-portal-api/requests"
	"gitlab.finema.co/finema/etda/web-portal-api/services"
	"net/http"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type MeController struct{}

func (a MeController) Profile(c core.IHTTPContext) error {
	user := c.Get(consts.ContextKeyUser).(*models.User)
	userSvc := services.NewUserService(c)
	item, ierr := userSvc.Find(user.ID)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, item)
}

func (a MeController) ChangePassword(c core.IHTTPContext) error {
	input := &requests.UserPasswordChange{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	user := c.Get(consts.ContextKeyUser).(*models.User)
	userSvc := services.NewUserService(c)
	ierr := userSvc.ChangePassword(&services.UserPasswordChangePayload{
		ID:              user.ID,
		CurrentPassword: utils.GetString(input.CurrentPassword),
		NewPassword:     utils.GetString(input.NewPassword),
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.NoContent(http.StatusNoContent)
}
