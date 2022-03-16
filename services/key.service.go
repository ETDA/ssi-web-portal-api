package services

import (
	"time"

	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/emsgs"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	"gitlab.finema.co/finema/etda/web-portal-api/views"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type KeyStorePayload struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
	KeyType    string `json:"key_type"`
}

type KeySignResponse struct {
	Signature string `json:"signature"`
	Message   string `json:"message"`
}

type IKeyService interface {
	Store(payload *KeyStorePayload) (*models.Key, core.IError)
	Generate() (*models.Key, core.IError)
	GenerateRSA() (*models.Key, core.IError)
	Sign(id string, message string) (*views.KeySign, core.IError)
	SignJSON(id string, message interface{}) (*views.KeySign, core.IError)
}

type keyService struct {
	ctx core.IContext
}

var keyRepositoryServiceConnectTimeout = 2 * time.Minute

func NewKeyService(ctx core.IContext) IKeyService {
	return &keyService{ctx: ctx}
}

func (s keyService) Store(payload *KeyStorePayload) (*models.Key, core.IError) {
	keyStoreResponse := &models.Key{}

	ierr := core.RequesterToStruct(keyStoreResponse, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Post("/key/store", payload, &core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVKeyServiceBaseURL),
			Timeout: &keyRepositoryServiceConnectTimeout,
		})
	})
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	return keyStoreResponse, nil
}

func (s keyService) Generate() (*models.Key, core.IError) {
	keyStoreResponse := &models.Key{}

	ierr := core.RequesterToStruct(keyStoreResponse, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Post("/key/generate", core.Map{}, &core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVKeyServiceBaseURL),
			Timeout: &keyRepositoryServiceConnectTimeout,
		})
	})
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	return keyStoreResponse, nil
}

func (s keyService) GenerateRSA() (*models.Key, core.IError) {
	keyStoreResponse := &models.Key{}

	ierr := core.RequesterToStruct(keyStoreResponse, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Post("/key/generate/rsa", core.Map{}, &core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVKeyServiceBaseURL),
			Timeout: &keyRepositoryServiceConnectTimeout,
		})
	})
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	return keyStoreResponse, nil
}

func (s keyService) Sign(id string, message string) (*views.KeySign, core.IError) {
	if id == "" {
		return nil, s.ctx.NewError(emsgs.KeyIDIsEmpty, emsgs.KeyIDIsEmpty)
	}
	if message == "" {
		return nil, s.ctx.NewError(emsgs.KeyMessageIsEmpty, emsgs.KeyMessageIsEmpty)
	}

	keySignResponse := &KeySignResponse{}
	ierr := core.RequesterToStruct(keySignResponse, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Post("/key/sign", core.Map{
			"id":      id,
			"message": message,
		}, &core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVKeyServiceBaseURL),
			Timeout: &keyRepositoryServiceConnectTimeout,
		})
	})
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	return views.NewKeySign(keySignResponse.Signature, keySignResponse.Message), nil
}

func (s keyService) SignJSON(id string, message interface{}) (*views.KeySign, core.IError) {
	if id == "" {
		return nil, s.ctx.NewError(emsgs.KeyIDIsEmpty, emsgs.KeyIDIsEmpty)
	}
	if message == "" {
		return nil, s.ctx.NewError(emsgs.KeyMessageIsEmpty, emsgs.KeyMessageIsEmpty)
	}

	messageBase64 := utils.Base64Encode(utils.JSONToString(message))
	keySignResponse := &KeySignResponse{}
	ierr := core.RequesterToStruct(keySignResponse, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Post("/key/sign", core.Map{
			"id":      id,
			"message": messageBase64,
		}, &core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVKeyServiceBaseURL),
			Timeout: &keyRepositoryServiceConnectTimeout,
		})
	})
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	return views.NewKeySign(keySignResponse.Signature, keySignResponse.Message), nil
}
