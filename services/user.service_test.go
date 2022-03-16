package services

import (
	"github.com/stretchr/testify/suite"
	core "ssi-gitlab.teda.th/ssi/core"
	"testing"
)

type UserServiceTestSuite struct {
	suite.Suite
	ctx *core.ContextMock
	s   IUserService
}

func TestUserServiceTestSuite(t *testing.T) {
	suite.Run(t, new(UserServiceTestSuite))
}

func (s *UserServiceTestSuite) SetupTest() {
	s.ctx = core.NewMockContext()
	s.ctx.On("DB").Return(s.ctx.MockDB.Gorm)

	s.s = NewUserService(s.ctx)
}
