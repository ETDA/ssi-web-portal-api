package web

import (
	"fmt"
	"net/http"

	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/helpers"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	"gitlab.finema.co/finema/etda/web-portal-api/requests"
	"gitlab.finema.co/finema/etda/web-portal-api/services"
	"gitlab.finema.co/finema/etda/web-portal-api/views"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type VCController struct{}

func (a VCController) VerifyVC(c core.IHTTPContext) error {

	input := &requests.VCVerify{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	vcService := services.NewVCService(c)
	res, ierr := vcService.VerifyVC(&services.VCVerifyVCPayload{Message: utils.GetString(input.JWT)})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, res)

}

func (a VCController) VerifyVP(c core.IHTTPContext) error {
	input := &requests.VPVerify{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	vcService := services.NewVCService(c)
	res, ierr := vcService.VerifyVP(&services.VCVerifyVPPayload{Message: utils.GetString(input.JWT)})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, res)
}

func (a VCController) Create(c core.IHTTPContext) error {
	input := &requests.VCSignRequest{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}
	vcService := services.NewVCService(c)
	view, ierr := vcService.CreateSignRequest(&services.VCSignRequestPayload{
		SchemaName:           utils.GetString(input.SchemaName),
		Signer:               utils.GetString(input.Signer),
		Holder:               utils.GetString(input.Holder),
		CredentialSubject:    input.CredentialSubject,
		CredentialSchemaID:   utils.GetString(input.CredentialSchema.ID),
		CredentialSchemaType: utils.GetString(input.CredentialSchema.Type),
		CreatorID:            c.Get(consts.ContextKeyUser).(*models.User).ID,
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, view)
}

func (a VCController) Find(c core.IHTTPContext) error {
	vcService := services.NewVCService(c)
	vc, ierr := vcService.Find(c.Param("id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	if vc.CID == "" {
		return c.JSON(http.StatusOK, views.NewVCItem(vc, nil))
	}
	status, ierr := vcService.GetVCStatus(vc.CID)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, views.NewVCItem(vc, status.Tags))
}

func (a VCController) FindByDID(c core.IHTTPContext) error {
	vcService := services.NewVCService(c)
	items, pageResponse, ierr := vcService.PaginationByDID(c.Param("did"), c.GetPageOptions())
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, core.NewPagination(items, pageResponse))
}
func (a VCController) Pagination(c core.IHTTPContext) error {
	vcService := services.NewVCService(c)
	items, pageResponse, ierr := vcService.Pagination(
		c.QueryParam("status"),
		c.QueryParam("start_date"),
		c.QueryParam("end_date"),
		c.GetPageOptionsWithOptions(&core.PageOptionsOptions{}))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, core.NewPagination(items, pageResponse))
}
func (a VCController) Revoke(c core.IHTTPContext) error {

	user := c.Get(consts.ContextKeyUser).(*models.User)
	vcService := services.NewVCService(c)

	ierr := vcService.RevokeByServer(user.Organization.ID, c.Param("id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.NoContent(http.StatusNoContent)
}

func (a VCController) CreateVCQR(c core.IHTTPContext) error {
	input := &requests.VCQRTokenRequest{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	vcService := services.NewVCService(c)
	qrToken, ierr := vcService.CreateQRToken(input.CIDs, utils.GetString(input.DIDAddress))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	qrCodeURL := fmt.Sprintf("%s/api/web/vcs/qr/%s", c.ENV().String(consts.ENVWebPortalBaseURL), qrToken.ID)
	return c.JSON(http.StatusOK, views.NewVCQRView(qrToken, qrCodeURL))
}

func (a VCController) VerifyVCQR(c core.IHTTPContext) error {
	qrToken := c.Get(consts.ContextKeyQRToken).(*models.VCQRToken)
	vcService := services.NewVCService(c)
	ierr := vcService.DeleteQRToken(qrToken.ID)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	orgSvc := services.NewOrganizationService(c)
	org, ierr := orgSvc.First()
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	vcs := &views.VCQRVerifyView{}
	vcs.CreatedAt = qrToken.CreatedAt
	vcs.SenderDID = org.Name
	cids := helpers.GetJSONValue(qrToken.CIDs, "cids").([]interface{})
	for _, cid := range cids {
		vc, ierr := vcService.FindByCID(cid.(string))
		if ierr != nil {
			return c.JSON(ierr.GetStatus(), ierr.JSON())
		}
		vcs.VCs = append(vcs.VCs, vc.JWT)
	}

	return c.JSON(http.StatusOK, vcs)
}

func (a VCController) Update(c core.IHTTPContext) error {
	input := &requests.VCUpdate{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}
	vcService := services.NewVCService(c)
	vc, ierr := vcService.Update(&services.VCUpdatePayload{
		ID:     c.Param("id"),
		Status: utils.GetString(input.Status),
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, vc)
}

func (a VCController) Approve(c core.IHTTPContext) error {
	input := &requests.VCApprove{}
	if err := c.BindWithValidateMessage(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	vcService := services.NewVCService(c)
	_, ierr := vcService.Approve(c.Param("id"), &services.VCApprovePayload{
		JWT: utils.GetString(input.JWT),
	})

	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, core.Map{
		"status": "success",
	})
}

func (a VCController) Reject(c core.IHTTPContext) error {
	input := &requests.VCReject{}
	if err := c.BindWithValidateMessage(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}
	vcService := services.NewVCService(c)
	_, ierr := vcService.Reject(c.Param("id"), &services.VCRejectPayload{
		RejectedReason: utils.GetString(input.Reason),
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, core.Map{
		"status": "success",
	})
}
