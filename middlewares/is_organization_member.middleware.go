package middlewares

import (
	"github.com/labstack/echo/v4"
	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/emsgs"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	"gitlab.finema.co/finema/etda/web-portal-api/services"
	core "ssi-gitlab.teda.th/ssi/core"
)

func IsOrganizationMember(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(core.IHTTPContext)
		user := cc.Get(consts.ContextKeyUser).(*models.User)
		organizationSvc := services.NewOrganizationService(cc)
		organization, ierr := organizationSvc.First()
		if ierr != nil {
			return c.JSON(ierr.GetStatus(), ierr.JSON())
		}
		orgID := organization.ID
		if user.Organization == nil {
			return c.JSON(emsgs.OrganizationPermissionDenied.GetStatus(), emsgs.OrganizationPermissionDenied.JSON())
		}

		if user.Organization.ID != orgID {
			return c.JSON(emsgs.OrganizationPermissionDenied.GetStatus(), emsgs.OrganizationPermissionDenied.JSON())
		}

		return next(c)
	}
}
