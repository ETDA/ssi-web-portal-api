package services

import (
	"github.com/stretchr/testify/mock"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
)

type MockOrganizationService struct {
	mock.Mock
}

func (m *MockOrganizationService) GenerateKey(id string) (*models.Organization, core.IError) {
	args := m.Called(id)
	return args.Get(0).(*models.Organization), core.MockIError(args, 1)
}

func (m *MockOrganizationService) StoreKey(id string, payload *KeyStorePayload) (*models.Organization, core.IError) {
	args := m.Called(id, payload)
	return args.Get(0).(*models.Organization), core.MockIError(args, 1)
}

func (m *MockOrganizationService) Pagination(pageOptions *core.PageOptions) ([]models.Organization, *core.PageResponse, core.IError) {
	args := m.Called(pageOptions)
	return args.Get(0).([]models.Organization), args.Get(1).(*core.PageResponse), core.MockIError(args, 2)
}

func (m *MockOrganizationService) Find(id string) (*models.Organization, core.IError) {
	args := m.Called(id)
	return args.Get(0).(*models.Organization), core.MockIError(args, 1)
}

func (m *MockOrganizationService) Create(payload *OrganizationCreatePayload) (*models.Organization, core.IError) {
	args := m.Called(payload)
	return args.Get(0).(*models.Organization), core.MockIError(args, 1)
}

func NewMockOrganizationService() *MockOrganizationService {
	return &MockOrganizationService{}
}
