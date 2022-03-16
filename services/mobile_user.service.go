package services

import (
	"errors"
	"fmt"
	"net/http"

	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	"gitlab.finema.co/finema/etda/web-portal-api/views"
	"gorm.io/gorm"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"ssi-gitlab.teda.th/ssi/core/utils"

	"gitlab.finema.co/finema/etda/web-portal-api/helpers"
	core "ssi-gitlab.teda.th/ssi/core"
)

type NotificationItem struct {
	Title    string `json:"title"`
	Body     string `json:"body"`
	ImageURL string `json:"image_url"`

	Category    string            `json:"category"`
	Icon        string            `json:"icon"`
	ClickAction string            `json:"click_action"`
	Sound       string            `json:"sound"`
	Priority    string            `json:"priority"` // one of "normal" or "high"
	Data        map[string]string `json:"data"`
}

type SendNotificationPayload struct {
	DidAddress    string             `json:"did_address"`
	Notifications []NotificationItem `json:"notifications"`
}

type IMobileUserService interface {
	Pagination(pageOptions *core.PageOptions) ([]views.MobileUser, *core.PageResponse, core.IError)
	SendNotification(payload *SendNotificationPayload) core.IError
	Find(id string) (*models.MobileUser, core.IError)
	FindByID(id string) (*views.MobileUserWithGroup, core.IError)
	FindByDID(did string) (*views.MobileUser, core.IError)
}

type mobileUserService struct {
	ctx             core.IContext
	walletConfigSvc IWalletConfigService
}

func NewMobileUserService(ctx core.IContext) IMobileUserService {
	return &mobileUserService{
		ctx:             ctx,
		walletConfigSvc: NewWalletConfigService(ctx),
	}
}

func (s mobileUserService) Pagination(pageOptions *core.PageOptions) ([]views.MobileUser, *core.PageResponse, core.IError) {
	config, ierr := s.walletConfigSvc.Get()
	if ierr != nil {
		return nil, nil, s.ctx.NewError(ierr, ierr)
	}
	users := make([]views.MobileUser, 0)
	queryParams := helpers.GetParamsFromPageOptions(pageOptions)
	pageResponse, ierr := core.RequesterToStructPagination(&users, pageOptions, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Get(fmt.Sprintf("/mobile/users"), &core.RequesterOptions{
			BaseURL: config.Endpoint,
			Params:  queryParams,
			Headers: http.Header{
				"Authorization": []string{fmt.Sprintf("%s %s", consts.AuthPrefix, config.AccessToken)},
			},
		})

	})
	if ierr != nil {
		return nil, nil, s.ctx.NewError(ierr, ierr)
	}
	for _, user := range users {
		_, err := s.Find(user.ID)
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

func (s mobileUserService) Find(id string) (*models.MobileUser, core.IError) {
	user := &models.MobileUser{}
	err := s.ctx.DB().First(user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, errmsgs.NotFoundCustomError("user"))
	}
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return user, nil

}

func (s mobileUserService) FindByID(id string) (*views.MobileUserWithGroup, core.IError) {
	config, ierr := s.walletConfigSvc.Get()
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	pageOptions := &core.PageOptions{}
	queryParams := helpers.GetParamsFromPageOptions(pageOptions)
	queryParams["ids"] = []string{id}
	users := make([]views.MobileUser, 0)
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
		return nil, s.ctx.NewError(ierr, ierr)
	}
	if pageResponse.Total == 0 {

		return nil, s.ctx.NewError(ierr, errmsgs.NotFoundCustomError("user"))
	}
	user := users[0]
	mobileUserGroupUsers := make([]models.MobileUserGroupUser, 0)
	err := s.ctx.DB().Find(&mobileUserGroupUsers, "mobile_user_id = ?", user.ID).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return views.NewMobileUserWithGroup(&user, make([]models.MobileUserGroup, 0)), nil
	}
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	var groupIDs []string
	for _, mobileUserGroupUser := range mobileUserGroupUsers {
		groupIDs = append(groupIDs, mobileUserGroupUser.MobileUserGroupID)
	}

	groups := make([]models.MobileUserGroup, 0)
	s.ctx.DB().Find(&groups, "id in (?)", groupIDs)
	return views.NewMobileUserWithGroup(&user, groups), nil

}

func (s mobileUserService) FindByDID(did string) (*views.MobileUser, core.IError) {
	config, ierr := s.walletConfigSvc.Get()
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	user := &views.MobileUser{}
	ierr = core.RequesterToStruct(user, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Get(fmt.Sprintf("/mobile/users/did_address/%s", did), &core.RequesterOptions{
			BaseURL: config.Endpoint,
			Headers: http.Header{
				"Authorization": []string{fmt.Sprintf("%s %s", consts.AuthPrefix, config.AccessToken)},
			},
		})
	})
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	return user, nil

}
func (s mobileUserService) SendNotification(payload *SendNotificationPayload) core.IError {
	config, ierr := s.walletConfigSvc.Get()
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}
	_, err := s.ctx.Requester().Post(fmt.Sprintf("/mobile/notification"), payload, &core.RequesterOptions{
		BaseURL: config.Endpoint,
		Headers: http.Header{
			"Authorization": []string{fmt.Sprintf("%s %s", consts.AuthPrefix, config.AccessToken)},
		},
	})
	if err != nil {
		return s.ctx.NewError(err, errmsgs.InternalServerError)
	}
	return nil
}
