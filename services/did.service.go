package services

import (
	"fmt"

	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/views"
	core "ssi-gitlab.teda.th/ssi/core"
)

type IDIDService interface {
	Create(keyID string, publicKey string) (*views.DIDDocument, core.IError)
	CreateWithRSAKey(keyID string, publicKey string) (*views.DIDDocument, core.IError)
	Find(didAddress string) (*views.DIDDocument, core.IError)
	GetNonce(did string) (string, core.IError)
}

type didService struct {
	ctx    core.IContext
	keySvc IKeyService
}

func NewDIDService(ctx core.IContext) IDIDService {
	return &didService{ctx: ctx, keySvc: NewKeyService(ctx)}
}

type GetNonceResponse struct {
	Nonce string `json:"nonce"`
}

func (s didService) Find(didAddress string) (*views.DIDDocument, core.IError) {
	did := &views.DIDDocument{}
	ierr := core.RequesterToStruct(did, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Get(fmt.Sprintf("/did/%s/document/latest", didAddress), &core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVDIDServiceBaseURL),
		})
	})

	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	return did, nil
}

func (s didService) GetNonce(did string) (string, core.IError) {
	nonce := &GetNonceResponse{}
	ierr := core.RequesterToStruct(nonce, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Get(fmt.Sprintf("/did/%s/nonce", did), &core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVDIDServiceBaseURL),
		})
	})

	if ierr != nil {
		return "", s.ctx.NewError(ierr, ierr)
	}

	return nonce.Nonce, nil
}

func (s didService) Create(keyID string, publicKey string) (*views.DIDDocument, core.IError) {
	message := core.Map{
		"operation":  consts.OperationDIDRegister,
		"public_key": publicKey,
		"key_type":   consts.KeyTypeSecp256r12019,
	}
	keySign, ierr := s.keySvc.SignJSON(keyID, message)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	did := &views.DIDDocument{}
	ierr = core.RequesterToStruct(did, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Post("/did", core.Map{
			"message": keySign.Message,
		}, &core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVDIDServiceBaseURL),
			Headers: map[string][]string{
				"x-signature": {keySign.Signature},
			},
		})
	})

	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	return did, nil
}

func (s didService) CreateWithRSAKey(keyID string, publicKey string) (*views.DIDDocument, core.IError) {
	message := core.Map{
		"operation":  consts.OperationDIDRegister,
		"public_key": publicKey,
		"key_type":   consts.KeyTypeRSA2018,
	}
	keySign, ierr := s.keySvc.SignJSON(keyID, message)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	did := &views.DIDDocument{}
	ierr = core.RequesterToStruct(did, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Post("/did", core.Map{
			"message": keySign.Message,
		}, &core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVDIDServiceBaseURL),
			Headers: map[string][]string{
				"x-signature": {keySign.Signature},
			},
		})
	})

	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	return did, nil
}
