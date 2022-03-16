package services

import (
	"errors"
	"fmt"
	"time"

	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	"gitlab.finema.co/finema/etda/web-portal-api/views"
	"gorm.io/gorm"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type IRequestedVPService interface {
	Find(id string) (*views.RequestedVPWithQRCode, core.IError)
	FindByQR(qrCodeID string) (*views.RequestedVP, core.IError)
	Pagination(status string, startDate string, endDate string, pageOptions *core.PageOptions) ([]views.RequestedVPWithQRCode, *core.PageResponse, core.IError)
	Create(creatorID string, payload *RequestedVPCreatePayload) (*views.RequestedVPWithQRCode, core.IError)
	Update(id string, payload *RequestedVPUpdatePayload) (*views.RequestedVPWithQRCode, core.IError)
	CancelList(payload *RequestVPCancelList) core.IError
	UpdateQRCode(id string) (*views.RequestedVPWithQRCode, core.IError)
}

type requestedVPService struct {
	ctx    core.IContext
	orgSvc IOrganizationService
}

func NewRequestedVPService(ctx core.IContext) IRequestedVPService {
	return &requestedVPService{
		ctx:    ctx,
		orgSvc: NewOrganizationService(ctx),
	}
}
func (s requestedVPService) Pagination(status string, startDate string, endDate string, pageOptions *core.PageOptions) ([]views.RequestedVPWithQRCode, *core.PageResponse, core.IError) {
	requestedVPs := make([]models.RequestedVP, 0)

	db := s.ctx.DB()
	var t time.Time
	const layout = "2006-01-02"
	if len(pageOptions.OrderBy) == 0 {
		pageOptions.OrderBy = []string{"created_at desc"}

	}
	if status != "" && status != "ALL" {
		db = db.Where("status = ?", status)
	}
	if startDate != "" {
		t, _ = time.Parse(layout, startDate)
		db = db.Where("created_at >= ?", t)
	}
	if endDate != "" {
		t, _ = time.Parse(layout, endDate)
		t.AddDate(0, 0, 1)
		db = db.Where("issuance_date <= ?", t)
	}
	db = core.SetSearchSimple(db, pageOptions.Q, []string{"name"})
	db = db.Preload("Creator")
	pageResponse, err := core.Paginate(db, &requestedVPs, pageOptions)
	if err != nil {
		return nil, nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	var requestedVPCounts []int64
	for _, requestedVP := range requestedVPs {
		var requestedVPCount int64
		err = s.ctx.DB().Find(&models.SubmittedVP{}, "requested_vp_id = ?", requestedVP.ID).Count(&requestedVPCount).Error
		requestedVPCounts = append(requestedVPCounts, requestedVPCount)
		if err != nil {
			return nil, nil, s.ctx.NewError(err, errmsgs.DBError)
		}
	}
	return views.NewRequestedVPList(requestedVPs, requestedVPCounts), pageResponse, nil
}
func (s requestedVPService) Create(creatorID string, payload *RequestedVPCreatePayload) (*views.RequestedVPWithQRCode, core.IError) {
	requestedVP := models.NewRequestedVP()
	requestedVP.CreatorID = creatorID
	requestedVP.Name = payload.Name
	requestedVP.QRCodeID = utils.NewSha256(requestedVP.ID)
	requestedVP.Schemacount = int64(len(payload.RequestSchemaList))
	err := s.ctx.DB().Create(requestedVP).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	for _, schema := range payload.RequestSchemaList {
		requestedVPSchema := models.NewRequestedVPSchemaType()
		requestedVPSchema.RequestedVPID = requestedVP.ID
		requestedVPSchema.SchemaType = schema.SchemaType
		requestedVPSchema.IsRequired = schema.IsRequired
		requestedVPSchema.Noted = schema.Noted
		err := s.ctx.DB().Create(requestedVPSchema).Error
		if err != nil {
			return nil, s.ctx.NewError(err, errmsgs.DBError)
		}
	}
	return s.Find(requestedVP.ID)
}

func (s requestedVPService) Update(id string, payload *RequestedVPUpdatePayload) (*views.RequestedVPWithQRCode, core.IError) {
	requestedVP := &models.RequestedVP{}
	err := s.ctx.DB().First(requestedVP, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, errmsgs.NotFoundCustomError("requestedVP"))
	}
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	requestedVP.Status = payload.Status
	requestedVP.UpdatedAt = utils.GetCurrentDateTime()
	err = s.ctx.DB().Updates(requestedVP).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	return s.Find(requestedVP.ID)
}

func (s requestedVPService) CancelList(payload *RequestVPCancelList) core.IError {
	for _, id := range payload.IDs {

		requestedVP := &models.RequestedVP{}
		err := s.ctx.DB().First(requestedVP, "id = ?", id).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return s.ctx.NewError(err, errmsgs.NotFoundCustomError("requestedVP"))
		}
		if err != nil {
			return s.ctx.NewError(err, errmsgs.DBError)
		}
		requestedVP.Status = consts.VPStatusCancel
		requestedVP.UpdatedAt = utils.GetCurrentDateTime()
		err = s.ctx.DB().Updates(requestedVP).Error
		if err != nil {
			return s.ctx.NewError(err, errmsgs.DBError)
		}
	}
	return nil
}
func (s requestedVPService) Find(id string) (*views.RequestedVPWithQRCode, core.IError) {

	requestedVP := &models.RequestedVP{}

	err := s.ctx.DB().Preload("Creator").First(requestedVP, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, errmsgs.NotFoundCustomError("requestedVP"))
	}

	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	requestedVPSchemaTypes := make([]models.RequestedVPSchemaType, 0)

	err = s.ctx.DB().Find(&requestedVPSchemaTypes, "requested_vp_id = ?", requestedVP.ID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, errmsgs.NotFoundCustomError("requestedVP"))
	}

	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	var submittedCount int64
	err = s.ctx.DB().Find(&models.SubmittedVP{}, "requested_vp_id = ?", requestedVP.ID).Count(&submittedCount).Error
	qrCodeURL := fmt.Sprintf("%s/api/web/requested-vps/qr/%s", s.ctx.ENV().String(consts.ENVWebPortalBaseURL), requestedVP.QRCodeID)
	return views.NewRequestedVPWithQRCode(requestedVP, requestedVPSchemaTypes, qrCodeURL, submittedCount), nil
}

func (s requestedVPService) FindByQR(qrCodeID string) (*views.RequestedVP, core.IError) {

	requestedVP := &models.RequestedVP{}

	err := s.ctx.DB().First(requestedVP, "qr_code_id = ?", qrCodeID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, errmsgs.NotFoundCustomError("requestedVP"))
	}

	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	requestedVPSchemaTypes := make([]models.RequestedVPSchemaType, 0)

	err = s.ctx.DB().Find(&requestedVPSchemaTypes, "requested_vp_id = ?", requestedVP.ID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, errmsgs.NotFoundCustomError("requestedVP"))
	}

	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	org, ierr := s.orgSvc.First()
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	requestedVPView := &views.RequestedVP{
		RequestedVP:           *requestedVP,
		RequestedVPSchemaType: requestedVPSchemaTypes,
		VerifierDID:           utils.GetString(org.DIDAddress),
		Verifier:              org.Name,
		Endpoint:              fmt.Sprintf("%s/api/web/requested-vps/%s/submit", s.ctx.ENV().String(consts.ENVWebPortalBaseURL), requestedVP.ID),
	}
	return requestedVPView, nil
}
func (s requestedVPService) UpdateQRCode(id string) (*views.RequestedVPWithQRCode, core.IError) {

	requestedVP := &models.RequestedVP{}

	err := s.ctx.DB().First(requestedVP, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, errmsgs.NotFoundCustomError("requestedVP"))
	}
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	requestedVP.QRCodeID = utils.NewSha256(utils.GetUUID())
	requestedVP.UpdatedAt = utils.GetCurrentDateTime()

	err = s.ctx.DB().Updates(requestedVP).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	return s.Find(requestedVP.ID)
}

type RequestedVPCreatePayload struct {
	Name              string                           `json:"name"`
	RequestSchemaList []RequestVPRequiredSchemaPayload `json:"schema_list"`
}
type RequestedVPUpdatePayload struct {
	Status string `json:"status"`
}
type RequestVPRequiredSchemaPayload struct {
	SchemaType string  `json:"schema_type"`
	IsRequired bool    `json:"is_required"`
	Noted      *string `json:"noted"`
}

type RequestVPCancelList struct {
	IDs []string `json:"ids"`
}
