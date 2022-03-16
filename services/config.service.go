package services

import (
	"errors"
	"fmt"
	"net/http"

	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/emsgs"

	"gitlab.finema.co/finema/etda/web-portal-api/models"
	"gorm.io/gorm"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type WalletConfigCreatePayload struct {
	Endpoint    string `json:"endpoint"`
	AccessToken string `json:"access_token"`
}
type IWalletConfigService interface {
	Get() (*models.WalletConfig, core.IError)
	Create(payload *WalletConfigCreatePayload) (*models.WalletConfig, core.IError)
	Find(id string) (*models.WalletConfig, core.IError)
	Delete(id string) core.IError
}
type walletConfigService struct {
	ctx core.IContext
}

func NewWalletConfigService(ctx core.IContext) IWalletConfigService {
	return &walletConfigService{ctx: ctx}
}

func (s walletConfigService) Get() (*models.WalletConfig, core.IError) {
	configs := &models.WalletConfig{}
	err := s.ctx.DB().First(configs).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, errmsgs.NotFoundCustomError("config"))
	}

	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	return configs, nil
}
func (s walletConfigService) Find(id string) (*models.WalletConfig, core.IError) {
	config := &models.WalletConfig{}
	err := s.ctx.DB().First(config, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, errmsgs.NotFoundCustomError("config"))
	}

	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	return config, nil
}
func (s walletConfigService) Create(payload *WalletConfigCreatePayload) (*models.WalletConfig, core.IError) {
	_, ierr := s.Get()
	if ierr == nil {
		return nil, s.ctx.NewError(ierr, emsgs.WalletConfigExists)
	}
	ierr = s.checkToken(payload.Endpoint, payload.AccessToken)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, emsgs.WalletConfigNotCorrect)
	}
	createdAt := utils.GetCurrentDateTime()
	id := utils.GetUUID()
	config := &models.WalletConfig{
		ID:          id,
		Endpoint:    payload.Endpoint,
		AccessToken: payload.AccessToken,
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}

	err := s.ctx.DB().Create(config).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	return s.Find(id)
}

type CheckStatusPayload struct {
	Status string `json:"status"`
}

func (s walletConfigService) checkToken(endpoint string, accessToken string) core.IError {

	status := &CheckStatusPayload{}
	ierr := core.RequesterToStruct(status, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Get("/mobile/status", &core.RequesterOptions{
			BaseURL: endpoint,
			Headers: http.Header{
				"Authorization": []string{fmt.Sprintf("%s %s", consts.AuthPrefix, accessToken)},
			},
		})
	})
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}
	return nil
}

func (s walletConfigService) Delete(id string) core.IError {

	config, ierr := s.Find(id)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}

	err := s.ctx.DB().Delete(config).Error
	if err != nil {
		return s.ctx.NewError(err, errmsgs.DBError)
	}
	err = s.ctx.DB().Where("1 = 1").Delete(&models.MobileUser{}).Error
	if err != nil {
		return s.ctx.NewError(err, errmsgs.DBError)
	}
	return nil
}
