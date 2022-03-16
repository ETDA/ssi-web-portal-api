package services

import (
	"errors"

	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/emsgs"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	"gitlab.finema.co/finema/etda/web-portal-api/views"
	"gorm.io/gorm"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type OrganizationCreatePayload struct {
	Name       string `json:"name"`
	JuristicID string `json:"juristic_id"`
}

type IOrganizationService interface {
	Pagination(pageOptions *core.PageOptions) ([]models.Organization, *core.PageResponse, core.IError)
	Find(id string) (*models.Organization, core.IError)
	Create(payload *OrganizationCreatePayload) (*models.Organization, core.IError)
	GenerateKey(id string) (*models.Organization, core.IError)
	GenerateWithRSAKey(id string) (*models.Organization, core.IError)
	StoreKey(id string, storeKeyPayload *KeyStorePayload) (*models.Organization, core.IError)
	First() (*models.Organization, core.IError)
}
type organizationService struct {
	ctx        core.IContext
	keyService IKeyService
	didService IDIDService
}

func NewOrganizationService(ctx core.IContext) IOrganizationService {
	return &organizationService{
		ctx:        ctx,
		keyService: NewKeyService(ctx),
		didService: NewDIDService(ctx)}
}

func (s organizationService) Pagination(pageOptions *core.PageOptions) ([]models.Organization, *core.PageResponse, core.IError) {
	items := make([]models.Organization, 0)

	db := s.ctx.DB()
	core.SetSearchSimple(db, pageOptions.Q, []string{"name", "juristic_id"})
	pageRes, err := core.Paginate(db, &items, pageOptions)
	if err != nil {
		return nil, nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return items, pageRes, nil
}

func (s organizationService) Find(id string) (*models.Organization, core.IError) {
	org := &models.Organization{}
	err := s.ctx.DB().First(org, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, errmsgs.NotFoundCustomError("organization"))
	}

	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return org, nil
}

func (s organizationService) Create(payload *OrganizationCreatePayload) (*models.Organization, core.IError) {
	org := models.NewOrganization()
	org.JuristicID = payload.JuristicID
	org.Name = payload.Name

	err := s.ctx.DB().Create(org).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return s.Find(org.ID)
}

func (s organizationService) GenerateKey(id string) (*models.Organization, core.IError) {
	org, ierr := s.Find(id)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	if org.EncryptedID != nil {
		return nil, s.ctx.NewError(emsgs.KeyAlreadyExists, emsgs.KeyAlreadyExists)
	}

	key, ierr := s.keyService.Generate()
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	did, ierr := s.didService.Create(key.ID, key.PublicKey)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	org.EncryptedID = &key.ID
	org.DIDAddress = &did.ID
	org.UpdatedAt = utils.GetCurrentDateTime()
	err := s.ctx.DB().Updates(org).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return s.Find(id)
}

func (s organizationService) GenerateWithRSAKey(id string) (*models.Organization, core.IError) {
	org, ierr := s.Find(id)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	if org.EncryptedID != nil {
		return nil, s.ctx.NewError(emsgs.KeyAlreadyExists, emsgs.KeyAlreadyExists)
	}

	key, ierr := s.keyService.GenerateRSA()
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	did, ierr := s.didService.CreateWithRSAKey(key.ID, key.PublicKey)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	org.EncryptedID = &key.ID
	org.DIDAddress = &did.ID
	org.UpdatedAt = utils.GetCurrentDateTime()
	err := s.ctx.DB().Updates(org).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return s.Find(id)
}
func (s organizationService) StoreKey(id string, payload *KeyStorePayload) (*models.Organization, core.IError) {
	org, ierr := s.Find(id)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	if org.EncryptedID != nil {
		return nil, s.ctx.NewError(emsgs.KeyAlreadyExists, emsgs.KeyAlreadyExists)
	}

	key, ierr := s.keyService.Store(payload)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	var did *views.DIDDocument
	if key.Type == string(consts.KeyTypeRSA) {
		did, ierr = s.didService.CreateWithRSAKey(key.ID, key.PublicKey)
		if ierr != nil {
			return nil, s.ctx.NewError(ierr, ierr)
		}
	}

	if key.Type == string(consts.KeyTypeECDSA) {
		did, ierr = s.didService.Create(key.ID, key.PublicKey)
		if ierr != nil {
			return nil, s.ctx.NewError(ierr, ierr)
		}
	}

	org.EncryptedID = &key.ID
	org.DIDAddress = &did.ID
	org.UpdatedAt = utils.GetCurrentDateTime()
	err := s.ctx.DB().Updates(org).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return s.Find(id)
}

func (s organizationService) First() (*models.Organization, core.IError) {
	org := &models.Organization{}
	err := s.ctx.DB().First(org).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, errmsgs.NotFoundCustomError("organization"))
	}

	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return org, nil
}
