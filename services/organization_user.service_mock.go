package services

import (
	"github.com/stretchr/testify/mock"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
)

type MockOrganizationUserService struct {
	mock.Mock
}

func (m *MockOrganizationUserService) Update(id string, orgID string, payload *OrganizationUserUpdatePayload) (*models.User, core.IError) {
	args := m.Called(id, orgID, payload)
	return args.Get(0).(*models.User), core.MockIError(args, 1)
}

func (m *MockOrganizationUserService) PaginationByOrganization(orgID string, pageOptions *core.PageOptions) ([]models.User, *core.PageResponse, core.IError) {
	args := m.Called(orgID, pageOptions)
	return args.Get(0).([]models.User), args.Get(1).(*core.PageResponse), core.MockIError(args, 2)
}

func (m *MockOrganizationUserService) FindByOrganization(id string, orgID string) (*models.User, core.IError) {
	args := m.Called(id, orgID)
	return args.Get(0).(*models.User), core.MockIError(args, 1)
}
