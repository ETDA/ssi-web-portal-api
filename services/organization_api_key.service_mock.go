package services

import (
	"github.com/stretchr/testify/mock"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
)

type MockOrganizationAPIKeyService struct {
	mock.Mock
}

func (m *MockOrganizationAPIKeyService) Pagination(pageOptions *core.PageOptions) ([]models.OrganizationAPIKey, *core.PageResponse, core.IError) {
	args := m.Called(pageOptions)
	return args.Get(0).([]models.OrganizationAPIKey), args.Get(1).(*core.PageResponse), core.MockIError(args, 2)
}

func (m *MockOrganizationAPIKeyService) Create(payload *OrganizationCreateAPIKeyPayload) (*models.OrganizationAPIKey, core.IError) {
	args := m.Called(payload)
	return args.Get(0).(*models.OrganizationAPIKey), core.MockIError(args, 1)
}

func (m *MockOrganizationAPIKeyService) Update(id string, payload *OrganizationUpdateAPIKeyPayload) (*models.OrganizationAPIKey, core.IError) {
	args := m.Called(id, payload)
	return args.Get(0).(*models.OrganizationAPIKey), core.MockIError(args, 1)
}

func (m *MockOrganizationAPIKeyService) Delete(id string) core.IError {
	args := m.Called(id)
	return core.MockIError(args, 0)
}

func (m *MockOrganizationAPIKeyService) Find(id string) (*models.OrganizationAPIKey, core.IError) {
	args := m.Called(id)
	return args.Get(0).(*models.OrganizationAPIKey), core.MockIError(args, 1)
}
