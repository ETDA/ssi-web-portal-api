//go:build e2e
// +build e2e

package services

import (
	"github.com/stretchr/testify/suite"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"testing"
)

type UserServiceE2ETestSuite struct {
	suite.Suite
	ctx                    core.IContext
	userSvc                IUserService
	testFindSuccess        core.Map
	testDeleteSuccess      core.Map
	testCreateTokenSuccess core.Map
}

func TestUserServiceE2ESuite(t *testing.T) {
	suite.Run(t, new(UserServiceE2ETestSuite))
}

func (s *UserServiceE2ETestSuite) AfterTest(suiteName, testName string) {
	if testName == "TestFind_Success" {
		if id, ok := s.testFindSuccess["organization_id"]; ok {
			s.NoError(s.ctx.DB().Unscoped().Where("id = ?", id).Delete(models.Organization{}).Error)
		}

		if id, ok := s.testFindSuccess["user_id"]; ok {
			s.NoError(s.ctx.DB().Unscoped().Where("id = ?", id).Delete(models.User{}).Error)
		}
	} else if testName == "TestDelete_Success" {
		if id, ok := s.testDeleteSuccess["organization_id"]; ok {
			s.NoError(s.ctx.DB().Unscoped().Where("id = ?", id).Delete(models.Organization{}).Error)
		}

		if id, ok := s.testDeleteSuccess["user_id"]; ok {
			s.NoError(s.ctx.DB().Unscoped().Where("id = ?", id).Delete(models.User{}).Error)
		}
	} else if testName == "TestCreateToken_Success" {
		if accessToken, ok := s.testDeleteSuccess["access_key"]; ok {
			s.NoError(s.ctx.DB().Unscoped().Where("access_key = ?", accessToken).Delete(models.OrganizationAPIKey{}).Error)
		}
	}
}

func (s *UserServiceE2ETestSuite) BeforeTest(suiteName, testName string) {
	env := core.NewENVPath("./../")

	mysql, _ := core.NewDatabase(env.Config()).Connect()
	s.ctx = core.NewContext(&core.ContextOptions{
		DB:  mysql,
		ENV: env,
	})

	s.userSvc = NewUserService(s.ctx)
}

func (s *UserServiceE2ETestSuite) TestFind_Success() {
	s.testFindSuccess = make(core.Map, 0)

	org := &models.Organization{
		ID:         utils.GetUUID(),
		JuristicID: utils.GetUUID(),
		Name:       "TestFind_Success",
		CreatedAt:  utils.GetCurrentDateTime(),
		UpdatedAt:  utils.GetCurrentDateTime(),
	}

	s.NoError(s.ctx.DB().Create(org).Error)

	s.testFindSuccess["organization_id"] = org.ID

	user := &models.User{
		ID:             utils.GetUUID(),
		OrganizationID: org.ID,
		Email:          "test_find_success@email.com",
		Password:       nil,
		Role:           "ADMIN",
		CreatedAt:      utils.GetCurrentDateTime(),
		UpdatedAt:      utils.GetCurrentDateTime(),
	}

	user.SetPassword("12345678")

	s.NoError(s.ctx.DB().Create(user).Error)

	s.testFindSuccess["user_id"] = user.ID

	findUser, ierr := s.userSvc.Find(user.ID)

	s.NoError(ierr)
	s.NotNil(findUser.Organization)
}

func (s *UserServiceE2ETestSuite) TestDelete_Success() {
	s.testDeleteSuccess = make(core.Map, 0)

	org := &models.Organization{
		ID:         utils.GetUUID(),
		JuristicID: utils.GetUUID(),
		Name:       "TestDelete_Success",
		CreatedAt:  utils.GetCurrentDateTime(),
		UpdatedAt:  utils.GetCurrentDateTime(),
	}

	s.NoError(s.ctx.DB().Create(org).Error)

	s.testDeleteSuccess["organization_id"] = org.ID

	user := &models.User{
		ID:             utils.GetUUID(),
		Email:          "test_delete_success@email.com",
		Password:       nil,
		OrganizationID: org.ID,
		Role:           "ADMIN",
		CreatedAt:      utils.GetCurrentDateTime(),
		UpdatedAt:      utils.GetCurrentDateTime(),
	}

	user.SetPassword("12345678")

	s.NoError(s.ctx.DB().Create(user).Error)

	s.testDeleteSuccess["user_id"] = user.ID

	ierr := s.userSvc.Delete(user.ID)
	s.NoError(ierr)

	findUser, ierr := s.userSvc.Find(user.ID)
	s.Error(ierr)
	s.True(errmsgs.IsNotFoundError(ierr))
	s.Nil(findUser)
}
