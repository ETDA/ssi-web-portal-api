package services

import (
	"fmt"
	"net/http"

	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/views"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
)

type ISchemaTokenService interface {
	Create(repositoryID string, payload *CreateSchemaTokenPayload) (*views.SchemaToken, core.IError)
	Update(repositoryID string, payload *UpdateSchemaTokenPayload) (*views.SchemaToken, core.IError)
	Find(repositoryID string, tokenID string) (*views.SchemaToken, core.IError)
	Pagination(repositoryID string, pageOptions *core.PageOptions) ([]views.SchemaToken, *core.PageResponse, core.IError)
	Delete(repositoryID string, tokenID string) core.IError
}
type schemaTokenService struct {
	ctx       core.IContext
	schemaSvc ISchemaService
}

func NewSchemaTokenService(ctx core.IContext) ISchemaTokenService {
	return &schemaTokenService{
		ctx:       ctx,
		schemaSvc: NewSchemaService(ctx),
	}
}

func (s schemaTokenService) Create(repositoryID string, payload *CreateSchemaTokenPayload) (*views.SchemaToken, core.IError) {
	schemaRepository, ierr := s.schemaSvc.FindRepository(repositoryID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	schemaToken := &views.SchemaToken{}
	ierr = core.RequesterToStruct(schemaToken, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Post(
			"/schemas/tokens",
			payload,
			&core.RequesterOptions{
				BaseURL: schemaRepository.Endpoint,
				Headers: http.Header{
					"Authorization": []string{fmt.Sprintf("%s %s", consts.AuthPrefix, schemaRepository.AccessToken)},
				},
			},
		)
	})

	return schemaToken, ierr
}
func (s schemaTokenService) Update(repositoryID string, payload *UpdateSchemaTokenPayload) (*views.SchemaToken, core.IError) {
	schemaRepository, ierr := s.schemaSvc.FindRepository(repositoryID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	schemaToken := &views.SchemaToken{}
	ierr = core.RequesterToStruct(schemaToken, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Put(
			fmt.Sprintf("/schemas/tokens/%s", payload.ID),
			payload,
			&core.RequesterOptions{
				BaseURL: schemaRepository.Endpoint,
				Headers: http.Header{
					"Authorization": []string{fmt.Sprintf("%s %s", consts.AuthPrefix, schemaRepository.AccessToken)},
				},
			},
		)
	})

	return schemaToken, ierr
}
func (s schemaTokenService) Find(repositoryID string, tokenID string) (*views.SchemaToken, core.IError) {
	schemaRepository, ierr := s.schemaSvc.FindRepository(repositoryID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	schemaToken := &views.SchemaToken{}
	ierr = core.RequesterToStruct(schemaToken, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Get(
			fmt.Sprintf("/schemas/tokens/%s", tokenID),
			&core.RequesterOptions{
				BaseURL: schemaRepository.Endpoint,
				Headers: http.Header{
					"Authorization": []string{fmt.Sprintf("%s %s", consts.AuthPrefix, schemaRepository.AccessToken)},
				},
			},
		)
	})
	return schemaToken, ierr
}
func (s schemaTokenService) Pagination(repositoryID string, pageOptions *core.PageOptions) ([]views.SchemaToken, *core.PageResponse, core.IError) {
	schemaRepository, ierr := s.schemaSvc.FindRepository(repositoryID)
	if ierr != nil {
		return nil, nil, s.ctx.NewError(ierr, ierr)
	}

	items := make([]views.SchemaToken, 0)
	pageResponse, ierr := core.RequesterToStructPagination(&items, pageOptions, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Get(
			"/schemas/tokens",
			&core.RequesterOptions{
				BaseURL: schemaRepository.Endpoint,
				Headers: http.Header{
					"Authorization": []string{fmt.Sprintf("%s %s", consts.AuthPrefix, schemaRepository.AccessToken)},
				},
			},
		)
	})
	return items, pageResponse, ierr

}
func (s schemaTokenService) Delete(repositoryID string, tokenID string) core.IError {
	schemaRepository, ierr := s.schemaSvc.FindRepository(repositoryID)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}
	// schemaToken := &views.SchemaToken{}
	_, err := s.ctx.Requester().Delete(
		fmt.Sprintf("/schemas/tokens/%s", tokenID),
		&core.RequesterOptions{
			BaseURL: schemaRepository.Endpoint,
			Headers: http.Header{
				"Authorization": []string{fmt.Sprintf("%s %s", consts.AuthPrefix, schemaRepository.AccessToken)},
			},
		},
	)

	if err != nil {
		return s.ctx.NewError(err, errmsgs.InternalServerError)
	}
	return nil
}

type CreateSchemaTokenPayload struct {
	Name string `json:"name"`
	Role string `json:"role"`
}

type UpdateSchemaTokenPayload struct {
	ID   string
	Name string `json:"name"`
	Role string `json:"role"`
}
