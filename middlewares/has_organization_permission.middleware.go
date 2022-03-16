package middlewares

import (
	"github.com/labstack/echo/v4"
	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/emsgs"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
)

func HasOrganizationPermission(permissions ...consts.UserRole) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := c.(core.IHTTPContext)
			user := cc.Get(consts.ContextKeyUser).(*models.User)

			isPass := false
			for _, permission := range permissions {
				if string(permission) == user.Role {
					isPass = true
					break
				}
			}

			if isPass == false {
				return c.JSON(emsgs.OrganizationPermissionDenied.GetStatus(), emsgs.OrganizationPermissionDenied.JSON())
			}

			return next(c)
		}
	}
}
