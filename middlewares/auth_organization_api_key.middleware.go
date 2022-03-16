package middlewares

import (
	"errors"
	"github.com/labstack/echo/v4"
	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/emsgs"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"gorm.io/gorm"
	"strings"
)

func AuthOrganizationAPIKey(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(core.IHTTPContext)
		authentication := strings.TrimSpace(cc.Request().Header.Get("Authorization"))
		if authentication == "" {
			return c.JSON(emsgs.AuthAPIKeyInvalid.GetStatus(), emsgs.AuthAPIKeyInvalid.JSON())
		}

		splittedAuthentication := strings.Split(authentication, " ")
		if len(splittedAuthentication) != 2 {
			return c.JSON(emsgs.AuthAPIKeyInvalid.GetStatus(), emsgs.AuthAPIKeyInvalid.JSON())
		}

		prefix := splittedAuthentication[0]
		accessKey := splittedAuthentication[1]
		if prefix != consts.AuthPrefix {
			return c.JSON(emsgs.AuthAPIKeyInvalid.GetStatus(), emsgs.AuthAPIKeyInvalid.JSON())
		}

		apiKey := &models.OrganizationAPIKey{}
		err := cc.DB().First(apiKey, "key = ?", accessKey).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(emsgs.AuthAPIKeyInvalid.GetStatus(), emsgs.AuthAPIKeyInvalid.JSON())
		}
		if err != nil {
			return c.JSON(errmsgs.DBError.GetStatus(), errmsgs.DBError.JSON())
		}

		return next(c)
	}
}
