package middlewares

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/emsgs"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	"gorm.io/gorm"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"strings"
)

func IsUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(core.IHTTPContext)
		authentication := strings.TrimSpace(cc.Request().Header.Get("Authorization"))
		if authentication == "" {
			return c.JSON(emsgs.AuthTokenRequired.GetStatus(), emsgs.AuthTokenRequired.JSON())
		}
		splittedAuthentication := strings.Split(authentication, " ")
		//if strings.splittedAuthentication

		if len(splittedAuthentication) != 2 {
			return c.JSON(emsgs.AuthTokenInvalid.GetStatus(), emsgs.AuthTokenInvalid.JSON())
		}
		prefix := splittedAuthentication[0]
		token := splittedAuthentication[1]
		if prefix != consts.AuthPrefix {
			return c.JSON(emsgs.AuthTokenInvalid.GetStatus(), emsgs.AuthTokenInvalid.JSON())
		}

		userAccessToken := &models.UserAccessToken{}
		err := cc.DB().First(userAccessToken, "token = ?", token).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(emsgs.AuthTokenInvalid.GetStatus(), emsgs.AuthTokenInvalid.JSON())
		}

		if err != nil {
			return c.JSON(errmsgs.DBError.GetStatus(), errmsgs.DBError.JSON())
		}

		user := &models.User{}
		err = cc.DB().Preload("Organization").First(user, "id = ? AND status = ?", userAccessToken.UserID, consts.UserStatusActive).Error
		if err != nil {
			return c.JSON(errmsgs.DBError.GetStatus(), errmsgs.DBError.JSON())
		}
		cc.Set(consts.ContextKeyUserToken, token)
		cc.Set(consts.ContextKeyUser, user)

		return next(c)
	}
}
