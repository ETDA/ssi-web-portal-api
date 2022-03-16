package middlewares

import (
	"github.com/labstack/echo/v4"
	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/emsgs"
	"gitlab.finema.co/finema/etda/web-portal-api/services"
	"log"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"strings"
)

func IsQRVerify(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(core.IHTTPContext)
		token := strings.TrimSpace(cc.Request().Header.Get("Authorization"))
		if token == "" {
			return c.JSON(emsgs.AuthTokenRequired.GetStatus(), emsgs.AuthTokenRequired.JSON())
		}

		log.Println(token)

		tokenID := cc.Param("token_id")
		vcService := services.NewVCService(cc)
		qrToken, ierr := vcService.FindQRToken(tokenID)

		if ierr != nil {
			return c.JSON(ierr.GetStatus(), ierr.JSON())
		}

		if qrToken.Token != token {
			// delete session
			ierr = vcService.DeleteQRToken(qrToken.ID)
			if ierr != nil {
				return c.JSON(ierr.GetStatus(), ierr.JSON())
			}

			// not found means this qr id with this token is not found.
			return c.JSON(errmsgs.NotFound.GetStatus(), errmsgs.NotFound.JSON())
		}

		cc.Set(consts.ContextKeyQRToken, qrToken)

		return next(c)
	}
}
