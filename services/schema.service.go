package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/emsgs"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	"gorm.io/gorm"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type ISchemaService interface {
	PaginationRepository(permission string, options *core.PageOptions) ([]models.VCSchemaConfig, *core.PageResponse, core.IError)
	FindRepository(id string) (*models.VCSchemaConfig, core.IError)
	CreateRepository(payloads []CreateSchemaRepositoryPayload) ([]models.VCSchemaConfig, core.IError)
	UpdateRepository(payload *UpdateSchemaRepositoryPayload) (*models.VCSchemaConfig, core.IError)
	DeleteRepository(id string) core.IError
	Pagination(repositoryID string, options *core.PageOptions) ([]models.VCSchema, *core.PageResponse, core.IError)
	Types(repositoryID string, schemaType string) (*SchemaTypes, core.IError)
	Create(repositoryID string, payload *CreateSchemaPayload) (*models.VCSchema, core.IError)
	Update(repositoryID string, schemaID string, payload *UpdateSchemaPayload) (*models.VCSchema, core.IError)
	Find(repositoryID string, schemaID string) (*models.VCSchema, core.IError)
	FindHistory(repositoryID string, schemaID string) ([]models.VCSchema, core.IError)
	FindByVersion(repositoryID string, schemaID string, version string) (*models.VCSchema, core.IError)
	FindReference(repositoryID string, schemaID string, version string, reference string) (*models.VCSchema, core.IError)
	CreateByUpload(repositoryID string, file core.IFile) ([]models.VCSchema, core.IError)
	UpdateByUpload(repositoryID string, schemaID string, file core.IFile) (*models.VCSchema, core.IError)
}

type CreateSchemaRepositoryPayload struct {
	Name     string
	Endpoint string
	Token    string
}

type UpdateSchemaRepositoryPayload struct {
	ID       string
	Name     string
	Endpoint string
	Token    string
}

type CreateSchemaPayload struct {
	SchemaName string           `json:"schema_name"`
	SchemaType string           `json:"schema_type"`
	SchemaBody *json.RawMessage `json:"schema_body"`
}

type UpdateSchemaPayload struct {
	SchemaName string           `json:"schema_name"`
	SchemaBody *json.RawMessage `json:"schema_body"`
	Version    string           `json:"version"`
}

type schemaService struct {
	ctx core.IContext
}

func NewSchemaService(ctx core.IContext) ISchemaService {
	return &schemaService{ctx: ctx}
}

func (s schemaService) PaginationRepository(permission string, options *core.PageOptions) ([]models.VCSchemaConfig, *core.PageResponse, core.IError) {
	items := make([]models.VCSchemaConfig, 0)
	db := s.ctx.DB()
	if permission == consts.SchemaRepositoryAdminRole ||
		permission == consts.SchemaRepositoryReadWriteRole ||
		permission == consts.SchemaRepositoryReadRole {
		db = db.Where("permission = ?", permission)
	}
	pageRes, err := core.Paginate(db, &items, options)
	if err != nil {
		return nil, nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return items, pageRes, nil
}

func (s *schemaService) FindRepository(id string) (*models.VCSchemaConfig, core.IError) {
	item := &models.VCSchemaConfig{}

	err := s.ctx.DB().Where("id = ?", id).First(item).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, emsgs.SchemaNotFound)
	}
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return item, nil
}

type vcSchemaRepositoryTokenModel struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	Token     string          `json:"token"`
	Role      string          `json:"role"`
	CreatedAt *time.Time      `json:"created_at"`
	UpdatedAt *time.Time      `json:"updated_at"`
	DeletedAt *gorm.DeletedAt `json:"deleted_at,omitempty"`
}

func (s schemaService) CreateRepository(payloads []CreateSchemaRepositoryPayload) ([]models.VCSchemaConfig, core.IError) {
	vcSchemaConfigs := make([]models.VCSchemaConfig, 0)
	for _, payload := range payloads {
		config := models.NewVCSchemaConfig(payload.Name, payload.Endpoint, payload.Token)
		token := &vcSchemaRepositoryTokenModel{}
		ierr := core.RequesterToStruct(token, func() (*core.RequestResponse, error) {
			return s.ctx.Requester().Get(
				"/schemas/tokens/me",
				&core.RequesterOptions{
					BaseURL: config.Endpoint,
					Headers: http.Header{
						"authorization": []string{fmt.Sprintf("%s %s", consts.AuthPrefix, config.AccessToken)},
					},
				},
			)
		})
		if ierr != nil {
			return nil, s.ctx.NewError(emsgs.SchemaRepositoryConectError(ierr), emsgs.SchemaRepositoryConectError(ierr))
		}

		config.Permission = token.Role

		err := s.ctx.DB().Create(config).Error
		if err != nil {
			return nil, s.ctx.NewError(err, errmsgs.DBError)
		}
		vcSchemaConfig, _ := s.FindRepository(config.ID)
		vcSchemaConfigs = append(vcSchemaConfigs, *vcSchemaConfig)
	}
	return vcSchemaConfigs, nil
}

func (s schemaService) UpdateRepository(payload *UpdateSchemaRepositoryPayload) (*models.VCSchemaConfig, core.IError) {
	config, ierr := s.FindRepository(payload.ID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	config.Name = payload.Name
	config.Endpoint = payload.Endpoint
	config.AccessToken = payload.Token

	token := &vcSchemaRepositoryTokenModel{}
	ierr = core.RequesterToStruct(token, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Get(
			"/schemas/tokens/me",
			&core.RequesterOptions{
				BaseURL: config.Endpoint,
				Headers: http.Header{
					"authorization": []string{fmt.Sprintf("%s %s", consts.AuthPrefix, config.AccessToken)},
				},
			},
		)
	})
	if ierr != nil {
		return nil, s.ctx.NewError(emsgs.SchemaRepositoryConectError(ierr), emsgs.SchemaRepositoryConectError(ierr))
	}
	config.Permission = token.Role
	config.UpdatedAt = utils.GetCurrentDateTime()
	err := s.ctx.DB().Updates(config).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return s.FindRepository(config.ID)
}

func (s schemaService) DeleteRepository(id string) core.IError {
	config, ierr := s.FindRepository(id)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}

	err := s.ctx.DB().Delete(config).Error
	if err != nil {
		return s.ctx.NewError(err, errmsgs.DBError)
	}

	return nil
}

func (s schemaService) Pagination(repositoryID string, options *core.PageOptions) ([]models.VCSchema, *core.PageResponse, core.IError) {
	var ierr core.IError
	schemaRepository := &models.VCSchemaConfig{}
	if repositoryID != "" {
		schemaRepository, ierr = s.FindRepository(repositoryID)
		if ierr != nil {
			return nil, nil, s.ctx.NewError(ierr, ierr)
		}
	}

	c := s.ctx.(core.IHTTPContext)
	items := make([]models.VCSchema, 0)
	pageRes, ierr := core.RequesterToStructPagination(&items, options, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Get(
			"/schemas",
			&core.RequesterOptions{
				BaseURL: schemaRepository.Endpoint,
				Headers: http.Header{
					"authorization": []string{fmt.Sprintf("%s %s", consts.AuthPrefix, schemaRepository.AccessToken)},
				},
				Params: c.QueryParams(),
			},
		)
	})
	for index, _ := range items {
		items[index].Permission = schemaRepository.Permission
	}
	return items, pageRes, ierr
}

type SchemaTypes struct {
	Types []string `json:"types"`
}

func (s schemaService) Types(repositoryID string, schemaType string) (*SchemaTypes, core.IError) {
	var ierr core.IError
	schemaRepository, ierr := s.FindRepository(repositoryID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	schemaTypes := &SchemaTypes{}
	ierr = core.RequesterToStruct(schemaTypes, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Get(
			fmt.Sprintf("/schemas/types?type=%s", schemaType),
			&core.RequesterOptions{
				BaseURL: schemaRepository.Endpoint,
				Headers: http.Header{
					"authorization": []string{fmt.Sprintf("%s %s", consts.AuthPrefix, schemaRepository.AccessToken)},
				},
			},
		)
	})

	return schemaTypes, ierr
}

func (s schemaService) Find(repositoryID string, schemaID string) (*models.VCSchema, core.IError) {
	schemaRepository, ierr := s.FindRepository(repositoryID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	schema := &models.VCSchema{}
	ierr = core.RequesterToStruct(schema, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Get(
			fmt.Sprintf("/schemas/%s", schemaID),
			&core.RequesterOptions{
				BaseURL: schemaRepository.Endpoint,
				Headers: http.Header{
					"authorization": []string{fmt.Sprintf("%s %s", consts.AuthPrefix, schemaRepository.AccessToken)},
				},
			},
		)
	})
	schema.Permission = schemaRepository.Permission

	return schema, ierr
}

func (s schemaService) FindHistory(repositoryID string, schemaID string) ([]models.VCSchema, core.IError) {
	schemaRepository, ierr := s.FindRepository(repositoryID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	schemas := make([]models.VCSchema, 0)
	ierr = core.RequesterToStruct(&schemas, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Get(
			fmt.Sprintf("/schemas/%s/history", schemaID),
			&core.RequesterOptions{
				BaseURL: schemaRepository.Endpoint,
				Headers: http.Header{
					"authorization": []string{fmt.Sprintf("%s %s", consts.AuthPrefix, schemaRepository.AccessToken)},
				},
			},
		)
	})
	for index, _ := range schemas {
		schemas[index].Permission = schemaRepository.Permission
	}

	return schemas, ierr
}

func (s schemaService) FindByVersion(repositoryID string, schemaID string, version string) (*models.VCSchema, core.IError) {
	schemaRepository, ierr := s.FindRepository(repositoryID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	schema := &models.VCSchema{}
	ierr = core.RequesterToStruct(schema, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Get(
			fmt.Sprintf("/schemas/%s/%s", schemaID, version),
			&core.RequesterOptions{
				BaseURL: schemaRepository.Endpoint,
				Headers: http.Header{
					"authorization": []string{fmt.Sprintf("%s %s", consts.AuthPrefix, schemaRepository.AccessToken)},
				},
			},
		)
	})
	schema.Permission = schemaRepository.Permission

	return schema, ierr
}

func (s schemaService) FindReference(repositoryID string, schemaID string, version string, reference string) (*models.VCSchema, core.IError) {
	schemaRepository, ierr := s.FindRepository(repositoryID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	schema := &models.VCSchema{}
	ierr = core.RequesterToStruct(schema, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Get(
			fmt.Sprintf("/schemas/%s/%s/%s", schemaID, version, reference),
			&core.RequesterOptions{
				BaseURL: schemaRepository.Endpoint,
				Headers: http.Header{
					"authorization": []string{fmt.Sprintf("%s %s", consts.AuthPrefix, schemaRepository.AccessToken)},
				},
			},
		)
	})
	schema.Permission = schemaRepository.Permission

	return schema, ierr
}

func (s schemaService) Create(repositoryID string, payload *CreateSchemaPayload) (*models.VCSchema, core.IError) {
	schemaRepository, ierr := s.FindRepository(repositoryID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	schema := &models.VCSchema{}
	ierr = core.RequesterToStruct(schema, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Post(
			"/schemas",
			payload,
			&core.RequesterOptions{
				BaseURL: schemaRepository.Endpoint,
				Headers: http.Header{
					"authorization": []string{fmt.Sprintf("%s %s", consts.AuthPrefix, schemaRepository.AccessToken)},
				},
			},
		)
	})
	schema.Permission = schemaRepository.Permission

	return schema, ierr
}

func (s schemaService) Update(repositoryID string, schemaID string, payload *UpdateSchemaPayload) (*models.VCSchema, core.IError) {
	schemaRepository, ierr := s.FindRepository(repositoryID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	schema := &models.VCSchema{}
	ierr = core.RequesterToStruct(schema, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Put(
			fmt.Sprintf("/schemas/%s", schemaID),
			payload,
			&core.RequesterOptions{
				BaseURL: schemaRepository.Endpoint,
				Headers: http.Header{
					"authorization": []string{fmt.Sprintf("%s %s", consts.AuthPrefix, schemaRepository.AccessToken)},
				},
			},
		)
	})
	schema.Permission = schemaRepository.Permission

	return schema, ierr
}

func (s schemaService) CreateByUpload(repositoryID string, file core.IFile) ([]models.VCSchema, core.IError) {
	schemaRepository, ierr := s.FindRepository(repositoryID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	schemas := make([]models.VCSchema, 0)
	ierr = core.RequesterToStruct(&schemas, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Post(
			fmt.Sprintf("/schemas/upload"),
			map[string]interface{}{
				"file": file,
			},
			&core.RequesterOptions{
				BaseURL: schemaRepository.Endpoint,
				Headers: http.Header{
					"authorization": []string{fmt.Sprintf("%s %s", consts.AuthPrefix, schemaRepository.AccessToken)},
				},
				IsMultipartForm: true,
			},
		)
	})
	for index, _ := range schemas {
		schemas[index].Permission = schemaRepository.Permission
	}

	return schemas, ierr
}

func (s schemaService) UpdateByUpload(repositoryID string, schemaID string, file core.IFile) (*models.VCSchema, core.IError) {
	schemaRepository, ierr := s.FindRepository(repositoryID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	schema := &models.VCSchema{}
	ierr = core.RequesterToStruct(schema, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Post(
			fmt.Sprintf("/schemas/%s/upload", schemaID),
			map[string]interface{}{
				"file": file,
			},
			&core.RequesterOptions{
				BaseURL: schemaRepository.Endpoint,
				Headers: http.Header{
					"authorization": []string{fmt.Sprintf("%s %s", consts.AuthPrefix, schemaRepository.AccessToken)},
				},
				IsMultipartForm: true,
			},
		)
	})
	schema.Permission = schemaRepository.Permission

	return schema, ierr
}
