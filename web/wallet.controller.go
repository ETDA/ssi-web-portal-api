package web

import (
	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	"gitlab.finema.co/finema/etda/web-portal-api/services"
	"net/http"
	"ssi-gitlab.teda.th/ssi/core"
)

type WalletController struct{}

func (a *WalletController) Summary(c core.IHTTPContext) error {
	organizationSvc := services.NewOrganizationService(c)
	keySvc := services.NewKeyService(c)
	s := services.NewWalletService(c, organizationSvc, keySvc)

	user := c.Get(consts.ContextKeyUser).(*models.User)

	summary, ierr := s.Summary(user.Organization.ID)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, summary)
}

func (a *WalletController) VCPagination(c core.IHTTPContext) error {
	organizationSvc := services.NewOrganizationService(c)
	keySvc := services.NewKeyService(c)
	s := services.NewWalletService(c, organizationSvc, keySvc)

	user := c.Get(consts.ContextKeyUser).(*models.User)

	vcs, pageResponse, ierr := s.VCPagination(user.Organization.ID, c.GetPageOptions(), &services.WalletVCPaginationOptions{})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, core.NewPagination(vcs, pageResponse))
}

func (a *WalletController) VCFind(c core.IHTTPContext) error {
	organizationSvc := services.NewOrganizationService(c)
	keySvc := services.NewKeyService(c)
	s := services.NewWalletService(c, organizationSvc, keySvc)

	user := c.Get(consts.ContextKeyUser).(*models.User)

	vc, ierr := s.VCFind(user.Organization.ID, c.Param("id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, vc)
}
