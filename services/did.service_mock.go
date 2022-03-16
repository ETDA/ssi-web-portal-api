package services

import (
	"github.com/stretchr/testify/mock"
	"gitlab.finema.co/finema/etda/web-portal-api/views"
	core "ssi-gitlab.teda.th/ssi/core"
)

type MockDIDService struct {
	mock.Mock
}

func (m *MockDIDService) GetNonce(did string) (string, core.IError) {
	args := m.Called(did)
	return args.String(0), core.MockIError(args, 1)
}

func (m *MockDIDService) Create(keyID string) (*views.DIDDocument, core.IError) {
	args := m.Called(keyID)
	return args.Get(0).(*views.DIDDocument), core.MockIError(args, 1)
}

func NewMockDIDService() *MockDIDService {
	return &MockDIDService{}
}
