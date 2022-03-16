package services

import (
	"github.com/stretchr/testify/mock"
	"gitlab.finema.co/finema/etda/web-portal-api/views"
	core "ssi-gitlab.teda.th/ssi/core"
)

type MockVCService struct {
	mock.Mock
}

func (m *MockVCService) VerifyVC(payload *VCVerifyVCPayload) (*views.VCVerify, core.IError) {
	args := m.Called(payload)
	return args.Get(0).(*views.VCVerify), core.MockIError(args, 1)
}

func (m *MockVCService) VerifyVP(payload *VCVerifyVPPayload) (*views.VPVerify, core.IError) {
	args := m.Called(payload)
	return args.Get(0).(*views.VPVerify), core.MockIError(args, 1)
}

func (m *MockVCService) Add(orgID string) (*views.VC, core.IError) {
	args := m.Called(orgID)
	return args.Get(0).(*views.VC), core.MockIError(args, 1)
}

func (m *MockVCService) Revoke(orgID string, cid string) core.IError {
	args := m.Called(orgID, cid)
	return core.MockIError(args, 0)
}

func NewMockVCService() *MockVCService {
	return &MockVCService{}
}
