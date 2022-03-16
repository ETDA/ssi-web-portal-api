package services

import (
	"errors"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"gorm.io/gorm"
)

type OrganizationUserUpdatePayload struct {
	Role string
}

type IOrganizationUserService interface {
	PaginationByOrganization(orgID string, pageOptions *core.PageOptions) ([]models.User, *core.PageResponse, core.IError)
	Find(id string, orgID string) (*models.User, core.IError)
	Update(id string, orgID string, payload *OrganizationUserUpdatePayload) (*models.User, core.IError)
	Delete(id string, orgID string) core.IError
}
type organizationUserService struct {
	ctx core.IContext
}

func NewOrganizationUserService(ctx core.IContext) IOrganizationUserService {
	return &organizationUserService{ctx: ctx}
}

func (s organizationUserService) Update(id string, orgID string, payload *OrganizationUserUpdatePayload) (*models.User, core.IError) {
	user, ierr := s.Find(id, orgID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	if payload.Role != "" {
		user.Role = payload.Role
	}

	err := s.ctx.DB().Updates(user).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return s.Find(id, orgID)
}

func (s organizationUserService) Find(id string, orgID string) (*models.User, core.IError) {
	user := &models.User{}
	err := s.ctx.DB().Where("organization_id = ? AND user_id = ?", orgID, id).Preload("Organization").First(user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, errmsgs.NotFoundCustomError("user"))
	}

	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return user, nil
}

func (s organizationUserService) PaginationByOrganization(orgID string, pageOptions *core.PageOptions) ([]models.User, *core.PageResponse, core.IError) {
	users := make([]models.User, 0)

	db := s.ctx.DB().Where("organization_id = ?", orgID).Preload("Organization")
	pageRes, err := core.Paginate(db, &users, pageOptions)
	if err != nil {
		return nil, nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return users, pageRes, nil
}

func (s organizationUserService) Delete(id string, orgID string) core.IError {
	user, ierr := s.Find(id, orgID)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}

	err := s.ctx.DB().Delete(&models.User{}, "user_id = ? AND organization_id = ?", id, orgID).Error
	if err != nil {
		return s.ctx.NewError(err, errmsgs.DBError)
	}

	err = s.ctx.DB().Delete(user).Error
	if err != nil {
		return s.ctx.NewError(err, errmsgs.DBError)
	}

	return nil
}
