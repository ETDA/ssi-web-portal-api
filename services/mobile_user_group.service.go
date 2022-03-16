package services

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"gitlab.finema.co/finema/etda/web-portal-api/helpers"
	"ssi-gitlab.teda.th/ssi/core/utils"

	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	"gitlab.finema.co/finema/etda/web-portal-api/views"
	"gorm.io/gorm"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
)

type IMobileUserGroupService interface {
	Create(payload *MobileUserGroupCreatePayload) (*views.MobileUserGroup, core.IError)
	Update(id string, payload *MobileUserGroupUpdatePayload) (*views.MobileUserGroup, core.IError)
	AddUser(groupIDs []string, userIDs []string) (*models.MobileUserGroup, core.IError)
	UpdateUser(groupIDs []string, userIDs string) core.IError
	RemoveUser(id string, userIDs []string) (*models.MobileUserGroup, core.IError)
	Find(id string) (*views.MobileUserGroup, core.IError)
	Get() ([]views.MobileUserGroup, core.IError)
	UserPagination(id string, pageOptions *core.PageOptions) ([]views.MobileUser, *core.PageResponse, core.IError)
	Delete(id string) core.IError
}

type mobileUserGroupService struct {
	ctx             core.IContext
	walletConfigSvc IWalletConfigService
	mobileUserSvc   IMobileUserService
}

func NewMobileUserGroupService(ctx core.IContext) IMobileUserGroupService {
	return &mobileUserGroupService{
		ctx:             ctx,
		walletConfigSvc: NewWalletConfigService(ctx),
		mobileUserSvc:   NewMobileUserService(ctx),
	}
}

func (s mobileUserGroupService) Find(id string) (*views.MobileUserGroup, core.IError) {
	group := &models.MobileUserGroup{}
	err := s.ctx.DB().First(group, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, errmsgs.NotFoundCustomError("group"))
	}
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	var userCount int64
	err = s.ctx.DB().Find(&models.MobileUserGroupUser{}, "mobile_user_group_id = ?", group.ID).Count(&userCount).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return views.NewMobileUserGroup(group, userCount), nil

}

func (s mobileUserGroupService) Create(payload *MobileUserGroupCreatePayload) (*views.MobileUserGroup, core.IError) {
	group := models.NewMobileUserGroup()
	group.Name = payload.Name
	err := s.ctx.DB().Create(group).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	return s.Find(group.ID)
}

func (s mobileUserGroupService) Update(id string, payload *MobileUserGroupUpdatePayload) (*views.MobileUserGroup, core.IError) {
	viewGroup, ierr := s.Find(id)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	group := &models.MobileUserGroup{}
	utils.Copy(group, viewGroup)
	if payload.Name != "" {
		group.Name = payload.Name
	}
	err := s.ctx.DB().Updates(group).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	return s.Find(group.ID)
}

func (s mobileUserGroupService) AddUser(GroupIDs []string, userIDs []string) (*models.MobileUserGroup, core.IError) {
	for _, groupID := range GroupIDs {
		for _, userID := range userIDs {
			var count int64
			err := s.ctx.DB().Find(&models.MobileUserGroupUser{}, "mobile_user_id = ? AND mobile_user_group_id = ?", userID, groupID).Count(&count).Error
			if count != 0 {
				continue
			}
			user := models.NewMobileUserGroupUser()
			user.MobileUserID = userID
			user.MobileUserGroupID = groupID
			err = s.ctx.DB().Create(user).Error
			if err != nil {
				return nil, s.ctx.NewError(err, errmsgs.DBError)
			}
		}
	}
	return nil, nil
}

func (s mobileUserGroupService) UpdateUser(groupIDs []string, userID string) core.IError {
	user, ierr := s.mobileUserSvc.Find(userID)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}

	err := s.ctx.DB().Delete(&models.MobileUserGroupUser{}, "mobile_user_id = ?", user.ID).Error
	if err != nil {
		return s.ctx.NewError(err, errmsgs.DBError)
	}
	_, err = s.AddUser(groupIDs, []string{user.ID})
	if err != nil {
		return s.ctx.NewError(err, errmsgs.DBError)
	}
	return nil
}
func (s mobileUserGroupService) RemoveUser(groupID string, userIDs []string) (*models.MobileUserGroup, core.IError) {
	for _, userID := range userIDs {
		err := s.ctx.DB().Delete(&models.MobileUserGroupUser{}, "mobile_user_group_id = ? AND mobile_user_id = ?", groupID, userID).Error
		if err != nil {
			return nil, s.ctx.NewError(err, errmsgs.DBError)
		}
	}
	return nil, nil
}
func (s mobileUserGroupService) Get() ([]views.MobileUserGroup, core.IError) {
	groups := make([]models.MobileUserGroup, 0)
	err := s.ctx.DB().Order("created_at ASC").Find(&groups).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, errmsgs.NotFoundCustomError("group"))
	}
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	var groupUserCounts []int64
	pageOptions := &core.PageOptions{}
	_, pageResponse, ierr := s.mobileUserSvc.Pagination(pageOptions)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	for _, group := range groups {
		var userCount int64
		err = s.ctx.DB().Find(&models.MobileUserGroupUser{}, "mobile_user_group_id = ?", group.ID).Count(&userCount).Error
		groupUserCounts = append(groupUserCounts, userCount)
		if err != nil {
			return nil, s.ctx.NewError(err, errmsgs.DBError)
		}
	}
	var distinctUserCount int64
	err = s.ctx.DB().Distinct("mobile_user_id").Find(&models.MobileUserGroupUser{}).Count(&distinctUserCount).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	groups = append([]models.MobileUserGroup{{
		ID:   consts.MobileUserALLGroup,
		Name: "ผู้ใช้งาน",
	},
	}, groups...)
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	groupUserCounts = append([]int64{pageResponse.Total}, groupUserCounts...)
	groupUserCounts = append(groupUserCounts, pageResponse.Total-distinctUserCount)
	groups = append(groups, models.MobileUserGroup{
		ID:   consts.MobileUserNoGroup,
		Name: "ผู้ใช้งานที่ยังไม่จัดกลุ่ม",
	})

	return views.NewMobileUserGroupList(groups, groupUserCounts), nil
}

func (s mobileUserGroupService) Delete(id string) core.IError {

	err := s.ctx.DB().Delete(&models.MobileUserGroup{}, "id = ?", id).Error
	if err != nil {
		return s.ctx.NewError(err, errmsgs.DBError)
	}
	return nil
}

func (s mobileUserGroupService) UserPagination(id string, pageOptions *core.PageOptions) ([]views.MobileUser, *core.PageResponse, core.IError) {
	config, ierr := s.walletConfigSvc.Get()

	if ierr != nil {
		return nil, nil, s.ctx.NewError(ierr, ierr)
	}
	var userIDs []string

	queryParams := helpers.GetParamsFromPageOptions(pageOptions)
	users := make([]views.MobileUser, 0)
	if id == consts.MobileUserALLGroup {
		return s.mobileUserSvc.Pagination(pageOptions)
	} else if id == consts.MobileUserNoGroup {
		err := s.ctx.DB().Model(&models.MobileUserGroupUser{}).Distinct("mobile_user_id").Pluck("mobile_user_id", &userIDs).Error
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, s.ctx.NewError(err, errmsgs.DBError)
		}
		if len(userIDs) >= 1 {
			queryParams["ignored_ids"] = []string{strings.Join(userIDs, ",")}
		}

	} else {
		group, ierr := s.Find(id)
		if ierr != nil {
			return nil, nil, s.ctx.NewError(ierr, ierr)
		}

		err := s.ctx.DB().Model(&models.MobileUserGroupUser{}).Where("mobile_user_group_id = ?", group.ID).Distinct("mobile_user_id").Pluck("mobile_user_id", &userIDs).Error
		if err != nil {
			return nil, nil, s.ctx.NewError(err, errmsgs.DBError)
		}

		queryParams["ids"] = []string{strings.Join(userIDs, ",")}
		if len(userIDs) == 0 {
			return users, &core.PageResponse{
				Total:   0,
				Limit:   pageOptions.Limit,
				Count:   0,
				Page:    1,
				Q:       pageOptions.Q,
				OrderBy: pageOptions.OrderBy,
			}, nil
		}
	}

	pageResponse, ierr := core.RequesterToStructPagination(&users, pageOptions, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Get(fmt.Sprintf("/mobile/users"), &core.RequesterOptions{
			BaseURL: config.Endpoint,
			Params:  queryParams,
			Headers: http.Header{
				"authorization": []string{fmt.Sprintf("%s %s", consts.AuthPrefix, config.AccessToken)},
			},
		})

	})
	if ierr != nil {
		return nil, nil, s.ctx.NewError(ierr, ierr)
	}

	for _, user := range users {
		_, err := s.mobileUserSvc.Find(user.ID)
		if err != nil {
			createdAt := utils.GetCurrentDateTime()
			mobileUser := &models.MobileUser{
				ID:        user.ID,
				CreatedAt: createdAt,
				UpdatedAt: createdAt,
			}
			ierr := s.ctx.DB().Create(mobileUser).Error
			if ierr != nil {
				return nil, nil, s.ctx.NewError(ierr, errmsgs.DBError)
			}
		}
	}
	return users, pageResponse, nil

}

type MobileUserGroupCreatePayload struct {
	Name string `json:"name"`
}
type MobileUserGroupUpdatePayload struct {
	Name string `json:"name"`
}
