package web

import (
	"net/http"

	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	"gitlab.finema.co/finema/etda/web-portal-api/requests"
	"gitlab.finema.co/finema/etda/web-portal-api/services"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type VPController struct{}

func (a VPController) Create(c core.IHTTPContext) error {
	//input := &requests.VPRequestCreate{}
	//if err := c.BindWithValidate(input); err != nil {
	//	return c.JSON(err.GetStatus(), err.JSON())
	//}
	panic("implement me")
}

func (a VPController) RequestedPagination(c core.IHTTPContext) error {
	requestedVPSvc := services.NewRequestedVPService(c)
	items, pageResponse, ierr := requestedVPSvc.Pagination(c.QueryParam("status"), c.QueryParam("start_date"), c.QueryParam("end_date"), c.GetPageOptions())
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON)
	}
	return c.JSON(http.StatusOK, core.NewPagination(items, pageResponse))
}
func (a VPController) RequestedCreate(c core.IHTTPContext) error {
	input := &requests.RequestedVPCreate{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}
	requestedVPSvc := services.NewRequestedVPService(c)
	payload := &services.RequestedVPCreatePayload{}
	utils.Copy(payload, input)
	utils.Copy(payload.RequestSchemaList, input.SchemaList)
	payload.RequestSchemaList = make([]services.RequestVPRequiredSchemaPayload, 0)
	for _, schema := range input.SchemaList {

		item := &services.RequestVPRequiredSchemaPayload{
			SchemaType: utils.GetString(schema.SchemaType),
			IsRequired: utils.GetBool(schema.IsRequired),
			Noted:      schema.Noted,
		}
		payload.RequestSchemaList = append(payload.RequestSchemaList, *item)
	}
	user := c.Get(consts.ContextKeyUser).(*models.User)
	requestedVP, ierr := requestedVPSvc.Create(user.ID, payload)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusCreated, requestedVP)
}
func (a VPController) RequestedUpdate(c core.IHTTPContext) error {
	input := &requests.RequestVPUpdate{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}
	requestedVPSvc := services.NewRequestedVPService(c)
	payload := &services.RequestedVPUpdatePayload{}
	utils.Copy(payload, input)
	requestedVP, ierr := requestedVPSvc.Update(c.Param("id"), payload)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, requestedVP)
}

func (a VPController) RequestUpdateList(c core.IHTTPContext) error {
	input := &requests.RequestVPCancelList{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}
	requestedVPSvc := services.NewRequestedVPService(c)
	payload := &services.RequestVPCancelList{}
	utils.Copy(payload, input)
	ierr := requestedVPSvc.CancelList(payload)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.NoContent(http.StatusNoContent)
}
func (a VPController) RequestedUpdateQRCode(c core.IHTTPContext) error {

	requestedVPSvc := services.NewRequestedVPService(c)
	requestedVP, ierr := requestedVPSvc.UpdateQRCode(c.Param("id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, requestedVP)
}
func (a VPController) RequestedFind(c core.IHTTPContext) error {
	requestedVPSvc := services.NewRequestedVPService(c)
	requestedVP, ierr := requestedVPSvc.Find(c.Param("id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, requestedVP)

}
func (a VPController) RequestedFindQR(c core.IHTTPContext) error {
	requestedVPSvc := services.NewRequestedVPService(c)
	requestedVP, ierr := requestedVPSvc.FindByQR(c.Param("id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	if requestedVP.Status != consts.VPStatusActive {
		err := errmsgs.NotFoundCustomError("Requested VP")
		return c.JSON(err.GetStatus(), err.JSON())
	}
	return c.JSON(http.StatusOK, requestedVP)
}
func (a VPController) SubmittedPagnination(c core.IHTTPContext) error {
	submittedVPSvc := services.NewSubmittedVPService(c)
	items, pageResponse, ierr := submittedVPSvc.PaginationByRequestedVP(
		c.Param("id"),
		c.QueryParam("start_date"),
		c.QueryParam("end_date"),
		c.GetPageOptions())
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON)
	}
	return c.JSON(http.StatusOK, core.NewPagination(items, pageResponse))
}

func (a VPController) SubmittedVCList(c core.IHTTPContext) error {
	submittedVPSvc := services.NewSubmittedVPService(c)
	submittedVPs, ierr := submittedVPSvc.GetVCList(c.Param("id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, core.NewPagination(submittedVPs, &core.PageResponse{
		Total:   int64(len(submittedVPs)),
		Limit:   1000,
		Count:   int64(len(submittedVPs)),
		Page:    1,
		Q:       "",
		OrderBy: nil,
	}))
}

func (a VPController) SubmittedFind(c core.IHTTPContext) error {
	submittedVPSvc := services.NewSubmittedVPService(c)
	submittedVP, ierr := submittedVPSvc.Find(c.Param("id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, submittedVP)
}

func (a VPController) SubmittedVCFind(c core.IHTTPContext) error {
	submittedVPSvc := services.NewSubmittedVPService(c)
	submittedVP, ierr := submittedVPSvc.GetVC(c.Param("id"), c.Param("vc_id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, submittedVP)
}

func (a VPController) TagStatus(c core.IHTTPContext) error {
	input := &requests.VPTagStatus{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}
	submittedVPSvc := services.NewSubmittedVPService(c)
	submittedVP, ierr := submittedVPSvc.TagStatus(c.Param("id"), input.Tags)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, submittedVP)
}
func (a VPController) SubmitVP(c core.IHTTPContext) error {
	jwt := c.Get(consts.ContextKeyJWTData).(*requests.JWTMessage)
	submittedVPSvc := services.NewSubmittedVPService(c)
	payload := &services.CreateSubmittedVPPayload{
		JWT:           c.Get(consts.ContextKeyJWT).(string),
		RequestedVPID: c.Param("id"),
		HolderDID:     jwt.Claims.Iss,
		DocumentCount: int64(len(jwt.Claims.VP.VerifiableCredential)),
		VCs:           jwt.Claims.VP.VerifiableCredential,
	}
	submittedVP, ierr := submittedVPSvc.Create(payload)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusCreated, submittedVP)
}
