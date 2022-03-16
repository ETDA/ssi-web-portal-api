package web

import (
	"net/http"

	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	"gitlab.finema.co/finema/etda/web-portal-api/requests"
	"gitlab.finema.co/finema/etda/web-portal-api/services"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type OrganizationController struct{}

func (a OrganizationController) Pagination(c core.IHTTPContext) error {
	orgSvc := services.NewOrganizationService(c)
	items, pageResponse, ierr := orgSvc.Pagination(c.GetPageOptions())
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, core.NewPagination(items, pageResponse))
}

func (a OrganizationController) Find(c core.IHTTPContext) error {
	orgSvc := services.NewOrganizationService(c)
	item, ierr := orgSvc.Find(c.Param("organization_id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, item)
}

func (a OrganizationController) Update(c core.IHTTPContext) error {
	return c.JSON(http.StatusOK, core.Map{
		"message": "Hello, I'm Home API",
	})
}

func (a OrganizationController) UploadCert(c core.IHTTPContext) error {
	return c.JSON(http.StatusOK, core.Map{
		"message": "Hello, I'm Home API",
	})
}

func (a OrganizationController) KeyGen(c core.IHTTPContext) error {
	orgSvc := services.NewOrganizationService(c)
	item, ierr := orgSvc.GenerateKey(c.Param("organization_id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusCreated, item)
}

func (a OrganizationController) UserPagination(c core.IHTTPContext) error {
	userSvc := services.NewOrganizationUserService(c)
	items, pageResponse, ierr := userSvc.PaginationByOrganization(c.Param("organization_id"), c.GetPageOptions())
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, core.NewPagination(items, pageResponse))
}

func (a OrganizationController) APIKeyCreate(c core.IHTTPContext) error {
	input := &requests.OrganizationAPIKeyCreate{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	service := services.NewOrganizationAPIKeyService(c)
	payload := &services.OrganizationCreateAPIKeyPayload{}
	_ = utils.Copy(payload, input)

	user := c.Get(consts.ContextKeyUser).(*models.User)
	payload.OrganizationID = user.OrganizationID
	apiKey, ierr := service.Create(payload)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusCreated, apiKey)
}

func (a OrganizationController) APIKeyPagination(c core.IHTTPContext) error {
	svc := services.NewOrganizationAPIKeyService(c)
	items, pageResponse, ierr := svc.Pagination(c.GetPageOptions())
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, core.NewPagination(items, pageResponse))
}

func (a OrganizationController) APIKeyFind(c core.IHTTPContext) error {
	svc := services.NewOrganizationAPIKeyService(c)
	item, ierr := svc.Find(c.Param("api_key_id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, item)
}

func (a OrganizationController) APIKeyUpdate(c core.IHTTPContext) error {
	input := &requests.OrganizationAPIKeyUpdate{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	svc := services.NewOrganizationAPIKeyService(c)
	item, ierr := svc.Update(c.Param("api_key_id"), &services.OrganizationUpdateAPIKeyPayload{
		Name: utils.GetString(input.Name),
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, item)
}

func (a OrganizationController) APIKeyDelete(c core.IHTTPContext) error {
	svc := services.NewOrganizationAPIKeyService(c)
	ierr := svc.Delete(c.Param("api_key_id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.NoContent(http.StatusNoContent)
}
