package services

import (
	"github.com/stretchr/testify/mock"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) PaginationByOrganization(orgID string, pageOptions *core.PageOptions) ([]models.User, *core.PageResponse, core.IError) {
	args := m.Called(orgID, pageOptions)
	return args.Get(0).([]models.User), args.Get(1).(*core.PageResponse), core.MockIError(args, 2)
}

func (m *MockUserService) FindByOrganization(id string, orgID string) (*models.User, core.IError) {
	args := m.Called(id, orgID)
	return args.Get(0).(*models.User), core.MockIError(args, 1)
}

func (m *MockUserService) Pagination(pageOptions *core.PageOptions) ([]models.User, *core.PageResponse, core.IError) {
	args := m.Called(pageOptions)
	return args.Get(0).([]models.User), args.Get(1).(*core.PageResponse), core.MockIError(args, 2)
}

func (m *MockUserService) Login(payload *UserLoginPayload) (*models.User, *models.UserAccessToken, core.IError) {
	args := m.Called(payload)
	return args.Get(0).(*models.User), args.Get(1).(*models.UserAccessToken), core.MockIError(args, 2)
}

func (m *MockUserService) Logout(token string) core.IError {
	args := m.Called(token)
	return core.MockIError(args, 0)
}

func (m *MockUserService) Update(id string, payload *UserUpdatePayload) (*models.User, core.IError) {
	args := m.Called(id, payload)
	return args.Get(0).(*models.User), core.MockIError(args, 1)
}

func (m *MockUserService) Verify(verifyToken string, password string) (*models.User, core.IError) {
	args := m.Called(verifyToken, password)
	return args.Get(0).(*models.User), core.MockIError(args, 1)
}

func (m *MockUserService) Register(payload *UserRegisterPayload) (*models.User, core.IError) {
	args := m.Called(payload)
	return args.Get(0).(*models.User), core.MockIError(args, 1)
}

func (m *MockUserService) Find(id string) (*models.User, core.IError) {
	args := m.Called(id)
	return args.Get(0).(*models.User), core.MockIError(args, 1)
}

func (m *MockUserService) ChangePassword(payload *UserPasswordChangePayload) core.IError {
	args := m.Called(payload)
	return core.MockIError(args, 0)
}
