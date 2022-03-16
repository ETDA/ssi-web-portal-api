package services

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/emsgs"
	"gitlab.finema.co/finema/etda/web-portal-api/helpers"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	"gitlab.finema.co/finema/etda/web-portal-api/views"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"gorm.io/gorm"
)

type ISubmittedVPService interface {
	Find(id string) (*models.SubmittedVP, core.IError)
	Create(payload *CreateSubmittedVPPayload) (*models.SubmittedVP, core.IError)
	PaginationByRequestedVP(requestedVPID string, startDate string, endDate string, pageOptions *core.PageOptions) ([]models.SubmittedVP, *core.PageResponse, core.IError)
	GetVCList(id string) ([]models.SubmittedVPVC, core.IError)
	GetVC(id string, vcID string) (*views.VCWithTag, core.IError)
	TagStatus(id string, tags []string) (*models.SubmittedVP, core.IError)
}

type submittedVPService struct {
	ctx            core.IContext
	mobileUserSvc  IMobileUserService
	schemaSvc      ISchemaService
	vcSvc          IVCService
	orgSvc         IOrganizationService
	keySvc         IKeyService
	didSvc         IDIDService
	requestedVPSvc IRequestedVPService
}

func NewSubmittedVPService(ctx core.IContext) ISubmittedVPService {
	return &submittedVPService{
		ctx:            ctx,
		mobileUserSvc:  NewMobileUserService(ctx),
		schemaSvc:      NewSchemaService(ctx),
		vcSvc:          NewVCService(ctx),
		orgSvc:         NewOrganizationService(ctx),
		keySvc:         NewKeyService(ctx),
		didSvc:         NewDIDService(ctx),
		requestedVPSvc: NewRequestedVPService(ctx),
	}
}

func (s submittedVPService) PaginationByRequestedVP(requestedVPID string, startDate string, endDate string, pageOptions *core.PageOptions) ([]models.SubmittedVP, *core.PageResponse, core.IError) {
	submittedVPs := make([]models.SubmittedVP, 0)

	db := s.ctx.DB().Where("requested_vp_id = ?", requestedVPID)
	var t time.Time
	const layout = "2006-01-02"
	if len(pageOptions.OrderBy) == 0 {
		pageOptions.OrderBy = []string{"created_at desc"}
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
	if len(pageOptions.OrderBy) == 0 {
		pageOptions.OrderBy = []string{"created_at desc"}
	}

	pageResponse, err := core.Paginate(db, &submittedVPs, pageOptions)
	if err != nil {
		return nil, nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	for index, item := range submittedVPs {
		user, ierr := s.mobileUserSvc.FindByDID(item.Holder)
		if ierr == nil {
			submittedVPs[index].Holder = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
		}
	}
	return submittedVPs, pageResponse, nil

}

func (s submittedVPService) GetVC(id string, vcID string) (*views.VCWithTag, core.IError) {
	vc := &models.SubmittedVPVC{}
	err := s.ctx.DB().First(&vc, "submitted_vp_id = ? AND id = ?", id, vcID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, errmsgs.NotFoundCustomError("vc"))
	}
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	vcStatus, ierr := s.vcSvc.GetVCStatus(vc.CID)

	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	return views.NewVCWithTag(vc, vcStatus.Tags), nil
}

func (s submittedVPService) GetVCList(id string) ([]models.SubmittedVPVC, core.IError) {
	vcs := make([]models.SubmittedVPVC, 0)
	err := s.ctx.DB().Find(&vcs, "submitted_vp_id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, errmsgs.NotFoundCustomError("vc"))
	}
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	return vcs, nil
}

func (s submittedVPService) Find(id string) (*models.SubmittedVP, core.IError) {
	submittedVP := &models.SubmittedVP{}
	err := s.ctx.DB().Preload("RequestedVP").Find(submittedVP, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, errmsgs.NotFoundCustomError("submitted VP"))
	}
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	user, ierr := s.mobileUserSvc.FindByDID(submittedVP.Holder)
	if ierr == nil {
		submittedVP.Holder = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	}
	return submittedVP, nil
}

func (s submittedVPService) Create(payload *CreateSubmittedVPPayload) (*models.SubmittedVP, core.IError) {
	requestedVP, ierr := s.requestedVPSvc.Find(payload.RequestedVPID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	if requestedVP.Status != consts.VPStatusActive {
		return nil, s.ctx.NewError(emsgs.RequestedVPSubmitError, emsgs.RequestedVPSubmitError)
	}

	verifyResult, ierr := s.vcSvc.VerifyVP(&VCVerifyVPPayload{
		Message: payload.JWT,
	})
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	for _, vc := range verifyResult.VC {
		if utils.GetString(vc.Status) != consts.VCStatusBlockchainActive {
			return nil, s.ctx.NewError(emsgs.VCNotActiveError, emsgs.VCNotActiveError)
		}
	}
	submittedVP := models.NewSubmittedVP()
	submittedVP.JWT = payload.JWT
	submittedVP.RequestedVPID = payload.RequestedVPID
	submittedVP.Holder = payload.HolderDID
	submittedVP.DocumentCount = payload.DocumentCount
	submittedVP.Verify = verifyResult.VerificationResult
	err := s.ctx.DB().Create(&submittedVP).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	for index, vcJWT := range payload.VCs {
		tokenM, _ := helpers.JWTVCDecodingT(vcJWT, []byte(""))
		if tokenM == nil || tokenM.Header == nil || tokenM.Claims == nil || tokenM.Signature == "" {
			return nil, s.ctx.NewError(emsgs.JWTInValid, emsgs.JWTInValid)
		}
		jwt := &JWTMessage{}
		_ = utils.MapToStruct(tokenM, jwt)
		schemaType := jwt.Claims.VC.Type[len(jwt.Claims.VC.Type)-1]
		issuanceDate := time.Unix(jwt.Claims.Iat, 0).UTC()
		vcSchema := &vcSchema{}
		ierr := core.RequesterToStruct(vcSchema, func() (*core.RequestResponse, error) {
			return s.ctx.Requester().Get(jwt.Claims.VC.CredentialSchema.ID, nil)
		})
		if ierr != nil {
			vcSchema.Name = ""
		}
		submittedVPVC := models.NewSubmittedVPVC()
		submittedVPVC.SubmittedVPID = submittedVP.ID
		submittedVPVC.CID = jwt.Claims.Jti
		submittedVPVC.SchemaName = vcSchema.Name
		submittedVPVC.SchemaType = schemaType
		submittedVPVC.IssuanceDate = &issuanceDate
		submittedVPVC.Issuer = jwt.Claims.Iss
		submittedVPVC.Holder = jwt.Claims.Sub
		submittedVPVC.JWT = vcJWT
		submittedVPVC.Status = consts.VCStatusActive
		submittedVPVC.Verify = verifyResult.VC[index].VerificationResult
		err = s.ctx.DB().Create(&submittedVPVC).Error
		if err != nil {
			return nil, s.ctx.NewError(err, errmsgs.DBError)
		}

	}
	return s.Find(submittedVP.ID)
}
func (s submittedVPService) TagStatus(id string, tags []string) (*models.SubmittedVP, core.IError) {
	org, ierr := s.orgSvc.First()
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	submittedVP, ierr := s.Find(id)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	var cids []string

	err := s.ctx.DB().Model(&models.SubmittedVPVC{}).Where("submitted_vp_id = ?", submittedVP.ID).Distinct("cid").Pluck("cid", &cids).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	nonce, ierr := s.didSvc.GetNonce(utils.GetString(org.DIDAddress))
	tagStatusRequest := core.Map{
		"operation":   consts.OperationVCTagStatus,
		"did_address": utils.GetString(org.DIDAddress),
		"cids":        cids,
		"tags":        tags,
		"nonce":       nonce,
	}
	keySign, ierr := s.keySvc.SignJSON(utils.GetString(org.EncryptedID), tagStatusRequest)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	did := &views.DIDDocument{}
	ierr = core.RequesterToStruct(did, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Post("/vc/status/tags", core.Map{
			"message": keySign.Message,
		}, &core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVVCStatusServiceBaseURL),
			Headers: map[string][]string{
				"x-signature": {keySign.Signature},
			},
		})
	})
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	submittedVP.Tags = strings.Join(tags, ",")
	submittedVP.UpdatedAt = utils.GetCurrentDateTime()
	err = s.ctx.DB().Updates(submittedVP).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	return s.Find(id)
}

type CreateSubmittedVPPayload struct {
	JWT           string
	RequestedVPID string
	HolderDID     string
	DocumentCount int64
	VCs           []string
}
type vcSchema struct {
	Name string `json:"name"`
}
type SubmittedVPVCList struct {
	VCs []string `json:"vcs"`
}
