package services

import (
	"github.com/stretchr/testify/mock"
	"gitlab.finema.co/finema/etda/web-portal-api/views"
	core "ssi-gitlab.teda.th/ssi/core"
)

type MockEKYCService struct {
	mock.Mock
}

func (m *MockDIDService) VerifyIDCard(payload *EKYCVerifyIDCardPayload) (*views.EKYCVerifyIDCard, core.IError) {
	args := m.Called(payload)
	return args.Get(0).(*views.EKYCVerifyIDCard), core.MockIError(args, 1)
}
