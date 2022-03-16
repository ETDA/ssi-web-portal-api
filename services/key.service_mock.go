package services

import (
	"github.com/stretchr/testify/mock"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	"gitlab.finema.co/finema/etda/web-portal-api/views"
	core "ssi-gitlab.teda.th/ssi/core"
)

type MockKeyService struct {
	mock.Mock
}

func (m *MockKeyService) Find(id string) (*models.Key, core.IError) {
	args := m.Called(id)
	return args.Get(0).(*models.Key), core.MockIError(args, 1)
}

func (m *MockKeyService) Store(payload *KeyStorePayload) (*models.Key, core.IError) {
	args := m.Called(payload)
	return args.Get(0).(*models.Key), core.MockIError(args, 1)
}

func (m *MockKeyService) Generate() (*models.Key, core.IError) {
	args := m.Called()
	return args.Get(0).(*models.Key), core.MockIError(args, 1)
}

func (m *MockKeyService) Sign(id string, message string) (*views.KeySign, core.IError) {
	args := m.Called(id, message)
	return args.Get(0).(*views.KeySign), core.MockIError(args, 1)
}

func (m *MockKeyService) SignJSON(id string, message interface{}) (*views.KeySign, core.IError) {
	args := m.Called(id, message)
	return args.Get(0).(*views.KeySign), core.MockIError(args, 1)
}

func NewMockKeyService() *MockKeyService {
	return &MockKeyService{}
}
