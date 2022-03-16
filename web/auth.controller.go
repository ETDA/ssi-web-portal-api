package web

import (
	"fmt"
	"net/http"

	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/requests"
	"gitlab.finema.co/finema/etda/web-portal-api/services"
	"gitlab.finema.co/finema/etda/web-portal-api/views"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type AuthController struct{}

func (a AuthController) Login(c core.IHTTPContext) error {
	input := &requests.UserLogin{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	service := services.NewUserService(c)
	browser, browserVersion := c.GetUserAgent().Browser()
	item, token, ierr := service.Login(&services.UserLoginPayload{
		Email:    utils.GetString(input.Email),
		Password: utils.GetString(input.Password),
		IP:       c.RealIP(),
		Device:   fmt.Sprintf("%s %s %s", c.GetUserAgent().OSInfo().FullName, browser, browserVersion),
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusCreated, &views.UserWithToken{
		User:  *item,
		Token: token.Token,
	})
}

func (a AuthController) Logout(c core.IHTTPContext) error {
	authSvc := services.NewUserService(c)
	ierr := authSvc.Logout(c.Get(consts.ContextKeyUserToken).(string))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusNoContent, nil)
}

func (a AuthController) ForgotPassword(c core.IHTTPContext) error {
	input := &requests.UserForgotPassword{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	//service := services.NewUserService(c)
	//ierr := service.ForgotPassword(utils.GetString(input.Email))
	//if ierr != nil {
	//	return c.JSON(ierr.GetStatus(), ierr.JSON())
	//}

	return c.NoContent(http.StatusNoContent)
}
