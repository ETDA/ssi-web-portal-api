package services

import (
	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/views"
	core "ssi-gitlab.teda.th/ssi/core"
)

type IVerifyService interface {
	VPVerify(jwt string) (*views.VPVerify, core.IError)
	VCVerify(jwt string) (*views.VCVerify, core.IError)
}

type verifyService struct {
	ctx core.IHTTPContext
}

func NewVerifyService(ctx core.IHTTPContext) IVerifyService {
	return &verifyService{ctx: ctx}
}

func (s verifyService) VPVerify(jwt string) (*views.VPVerify, core.IError) {
	payload := core.Map{
		"message": jwt,
	}
	res := &views.VPVerify{}
	ierr := core.RequesterToStruct(res, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Post("/vp/verify", payload, &core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVVerifyServiceBaseURL),
		})
	})
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	return res, nil
}

func (s verifyService) VCVerify(jwt string) (*views.VCVerify, core.IError) {
	payload := core.Map{
		"message": jwt,
	}
	res := &views.VCVerify{}
	ierr := core.RequesterToStruct(res, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Post("/vc/verify", payload, &core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVVerifyServiceBaseURL),
		})
	})
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	return res, nil
}
