package services

import (
	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/views"
	core "ssi-gitlab.teda.th/ssi/core"
)

type EKYCVerifyIDCardPayload struct {
	CardID    string `json:"card_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	LaserID   string `json:"laser_id"`
	Birthdate string `json:"birthdate"`
}

type IEKYCService interface {
	VerifyIDCard(payload *EKYCVerifyIDCardPayload) (*views.EKYCVerifyIDCard, core.IError)
}
type ekycService struct {
	ctx core.IContext
}

func NewEKYCService(ctx core.IContext) IEKYCService {
	return &ekycService{ctx: ctx}
}

func (s ekycService) VerifyIDCard(payload *EKYCVerifyIDCardPayload) (*views.EKYCVerifyIDCard, core.IError) {
	res := &views.EKYCVerifyIDCard{}
	ierr := core.RequesterToStruct(res, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Post("/id-card/verify", payload, &core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVVCIDProofingServiceBaseURL),
		})
	})
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	return res, nil
}
