package services

import (
	"fmt"
	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/helpers"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	"gitlab.finema.co/finema/etda/web-portal-api/views"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type WalletVCPaginationOptions struct {
	Role string
}

type IWalletService interface {
	Summary(orgID string) (*views.WalletSummary, core.IError)
	VCPagination(orgID string, pageOptions *core.PageOptions, options *WalletVCPaginationOptions) ([]models.VC, *core.PageResponse, core.IError)
	VCFind(orgID string, cid string) (*models.VC, core.IError)
}

type walletService struct {
	ctx             core.IContext
	organizationSvc IOrganizationService
	keySvc          IKeyService
}

func NewWalletService(ctx core.IContext, organizationSvc IOrganizationService, keySvc IKeyService) IWalletService {
	return &walletService{
		ctx:             ctx,
		organizationSvc: organizationSvc,
		keySvc:          keySvc,
	}
}

func (s walletService) Summary(orgID string) (*views.WalletSummary, core.IError) {
	org, ierr := s.organizationSvc.Find(orgID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	keySign, ierr := s.keySvc.Sign(utils.GetString(org.EncryptedID), utils.GetString(org.DIDAddress))
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	summary := &views.WalletSummary{}
	ierr = core.RequesterToStruct(summary, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Get(fmt.Sprintf("/wallet/%s", keySign.Message), &core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVWalletServiceBaseURL),
			Headers: map[string][]string{
				"x-signature": {keySign.Signature},
			},
		})
	})
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	return summary, nil
}

func (s walletService) VCPagination(orgID string, pageOptions *core.PageOptions, options *WalletVCPaginationOptions) ([]models.VC, *core.PageResponse, core.IError) {
	org, ierr := s.organizationSvc.Find(orgID)
	if ierr != nil {
		return nil, nil, s.ctx.NewError(ierr, ierr)
	}

	keySign, ierr := s.keySvc.Sign(utils.GetString(org.EncryptedID), utils.GetString(org.DIDAddress))
	if ierr != nil {
		return nil, nil, s.ctx.NewError(ierr, ierr)
	}

	items := make([]models.VC, 0)
	pagination, ierr := core.RequesterToStructPagination(&items, pageOptions, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Get(fmt.Sprintf("/wallet/%s/vcs", keySign.Message), &core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVWalletServiceBaseURL),
			Headers: map[string][]string{
				"x-signature": {keySign.Signature},
			},
			Params: helpers.GetParamsFromPageOptions(pageOptions),
		})
	})
	if ierr != nil {
		return nil, nil, s.ctx.NewError(ierr, ierr)
	}
	pagination.Count = int64(len(items))

	return items, pagination, nil
}

func (s walletService) VCFind(orgID string, cid string) (*models.VC, core.IError) {
	org, ierr := s.organizationSvc.Find(orgID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	keySign, ierr := s.keySvc.Sign(utils.GetString(org.EncryptedID), utils.GetString(org.DIDAddress))
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	vc := &models.VC{}
	ierr = core.RequesterToStruct(vc, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Get(fmt.Sprintf("/wallet/%s/vcs/%s", keySign.Message, cid), &core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVWalletServiceBaseURL),
			Headers: map[string][]string{
				"x-signature": {keySign.Signature},
			},
		})
	})
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	return vc, nil
}
