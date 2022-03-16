package services

import (
	"github.com/stretchr/testify/mock"
	core "ssi-gitlab.teda.th/ssi/core"
)

type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) WalletSummary(did string) (interface{}, core.IError) {
	args := m.Called(did)
	return args.Get(0), core.MockIError(args, 1)
}

func (m *MockWalletService) WalletVCPagination(did string, pageOptions *core.PageOptions, options *WalletVCPaginationOptions) (interface{}, core.IError) {
	args := m.Called(did, pageOptions, options)
	return args.Get(0), core.MockIError(args, 1)
}

func (m *MockWalletService) WalletVCFind(did string, cid string) (interface{}, core.IError) {
	args := m.Called(did, cid)
	return args.Get(0), core.MockIError(args, 1)
}

func NewMockWalletService() *MockWalletService {
	return &MockWalletService{}
}
