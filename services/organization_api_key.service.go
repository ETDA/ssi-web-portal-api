package services

import (
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type OrganizationCreateAPIKeyPayload struct {
	Name           string `json:"name"`
	Read           bool   `json:"read"`
	Write          bool   `json:"write"`
	OrganizationID string `json:"organization_id"`
}

type OrganizationUpdateAPIKeyPayload struct {
	Name string `json:"name"`
}
type IOrganizationAPIKeyService interface {
	Pagination(pageOptions *core.PageOptions) ([]models.OrganizationAPIKey, *core.PageResponse, core.IError)
	Create(payload *OrganizationCreateAPIKeyPayload) (*models.OrganizationAPIKey, core.IError)
	Update(id string, payload *OrganizationUpdateAPIKeyPayload) (*models.OrganizationAPIKey, core.IError)
	Delete(id string) core.IError
	Find(id string) (*models.OrganizationAPIKey, core.IError)
}
type organizationAPIKeyService struct {
	ctx core.IContext
}

func NewOrganizationAPIKeyService(ctx core.IContext) IOrganizationAPIKeyService {
	return &organizationAPIKeyService{ctx: ctx}
}

func (s organizationAPIKeyService) Pagination(pageOptions *core.PageOptions) ([]models.OrganizationAPIKey, *core.PageResponse, core.IError) {
	items := make([]models.OrganizationAPIKey, 0)
	db := s.ctx.DB()
	core.SetSearchSimple(db, pageOptions.Q, []string{"name"})
	pageRes, err := core.Paginate(db, &items, pageOptions)
	if err != nil {
		return nil, nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return items, pageRes, nil
}

func (s organizationAPIKeyService) Create(payload *OrganizationCreateAPIKeyPayload) (*models.OrganizationAPIKey, core.IError) {
	apiKey := models.NewAPIKey(payload.Name)
	apiKey.Read = payload.Read
	apiKey.Write = payload.Write
	apiKey.OrganizationID = payload.OrganizationID

	err := s.ctx.DB().Create(&apiKey).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return apiKey, nil
}

func (s organizationAPIKeyService) Update(id string, payload *OrganizationUpdateAPIKeyPayload) (*models.OrganizationAPIKey, core.IError) {
	apiKey, ierr := s.Find(id)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	apiKey.Name = payload.Name
	apiKey.UpdatedAt = utils.GetCurrentDateTime()
	err := s.ctx.DB().Updates(apiKey).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return apiKey, nil
}

func (s organizationAPIKeyService) Delete(id string) core.IError {
	apiKey, ierr := s.Find(id)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}

	err := s.ctx.DB().Delete(apiKey).Error
	if err != nil {
		return s.ctx.NewError(err, errmsgs.DBError)
	}

	return nil
}

func (s organizationAPIKeyService) Find(id string) (*models.OrganizationAPIKey, core.IError) {
	apiKey := &models.OrganizationAPIKey{}
	err := s.ctx.DB().First(apiKey, "id = ?", id).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return apiKey, nil
}
