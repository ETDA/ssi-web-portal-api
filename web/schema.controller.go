package web

import (
	"net/http"

	"gitlab.finema.co/finema/etda/web-portal-api/emsgs"
	"gitlab.finema.co/finema/etda/web-portal-api/helpers"

	"gitlab.finema.co/finema/etda/web-portal-api/requests"
	"gitlab.finema.co/finema/etda/web-portal-api/services"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type SchemaController struct{}

func (sc SchemaController) Pagination(c core.IHTTPContext) error {
	service := services.NewSchemaService(c)
	schemas, pageResponse, ierr := service.Pagination(c.Param("repository_id"), c.GetPageOptionsWithOptions(&core.PageOptionsOptions{}))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, core.NewPagination(schemas, pageResponse))
}

func (sc SchemaController) Types(c core.IHTTPContext) error {
	service := services.NewSchemaService(c)
	schemaTypes, ierr := service.Types(c.Param("repository_id"), c.QueryParam("q"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, schemaTypes)
}

func (sc SchemaController) Find(c core.IHTTPContext) error {
	service := services.NewSchemaService(c)
	schema, ierr := service.Find(c.Param("repository_id"), c.Param("schema_id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, schema)
}

func (sc SchemaController) FindHistory(c core.IHTTPContext) error {
	service := services.NewSchemaService(c)
	schema, ierr := service.FindHistory(c.Param("repository_id"), c.Param("schema_id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, schema)
}

func (sc SchemaController) FindByVersion(c core.IHTTPContext) error {
	service := services.NewSchemaService(c)
	schema, ierr := service.FindByVersion(c.Param("repository_id"), c.Param("schema_id"), c.Param("version"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, schema)
}

func (sc SchemaController) FindSchemaInstance(c core.IHTTPContext) error {
	service := services.NewSchemaService(c)
	schema, ierr := service.FindByVersion(c.Param("repository_id"), c.Param("schema_id"), c.Param("version"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, schema.SchemaBody)
}

func (sc SchemaController) FindSchemaReference(c core.IHTTPContext) error {
	service := services.NewSchemaService(c)
	schema, ierr := service.FindReference(c.Param("repository_id"), c.Param("schema_id"), c.Param("version"), c.Param("reference"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, schema)
}

func (sc SchemaController) Create(c core.IHTTPContext) error {
	input := &requests.SchemaCreate{}
	ierr := c.BindWithValidate(input)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	service := services.NewSchemaService(c)
	schema, ierr := service.Create(c.Param("repository_id"), &services.CreateSchemaPayload{
		SchemaName: utils.GetString(input.SchemaName),
		SchemaType: utils.GetString(input.SchemaType),
		SchemaBody: input.SchemaBody,
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusCreated, schema)
}

func (sc SchemaController) Update(c core.IHTTPContext) error {
	input := &requests.SchemaUpdate{}
	ierr := c.BindWithValidate(input)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	service := services.NewSchemaService(c)
	schema, ierr := service.Update(c.Param("repository_id"), c.Param("schema_id"), &services.UpdateSchemaPayload{
		SchemaName: utils.GetString(input.SchemaName),
		SchemaBody: input.SchemaBody,
		Version:    utils.GetString(input.Version),
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, schema)
}

func (sc SchemaController) CreateByUpload(c core.IHTTPContext) error {
	formFile, err := c.FormFile("file")
	if err != nil {
		return c.JSON(emsgs.RequestFormError(err).GetStatus(), emsgs.RequestFormError(err).JSON())
	}

	s := services.NewSchemaService(c)

	schemaFile, err := helpers.MultiPartFileToIFileHeader(formFile)
	if err != nil {
		return c.JSON(emsgs.OpenFileError(err).GetStatus(), emsgs.OpenFileError(err).JSON())
	}

	schema, ierr := s.CreateByUpload(c.Param("repository_id"), schemaFile)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusCreated, schema)
}

func (sc SchemaController) UploadByUpload(c core.IHTTPContext) error {
	formFile, err := c.FormFile("file")
	if err != nil {
		return c.JSON(emsgs.RequestFormError(err).GetStatus(), emsgs.RequestFormError(err).JSON())
	}

	s := services.NewSchemaService(c)

	schemaFile, err := helpers.MultiPartFileToIFileHeader(formFile)
	if err != nil {
		return c.JSON(emsgs.OpenFileError(err).GetStatus(), emsgs.OpenFileError(err).JSON())
	}

	schema, ierr := s.UpdateByUpload(c.Param("repository_id"), c.Param("schema_id"), schemaFile)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusOK, schema)
}
func (sc SchemaController) TokenPagination(c core.IHTTPContext) error {
	schemaTokenSvc := services.NewSchemaTokenService(c)
	items, pageResponse, ierr := schemaTokenSvc.Pagination(c.Param("repository_id"), c.GetPageOptions())
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, core.NewPagination(items, pageResponse))
}
func (sc SchemaController) TokenFind(c core.IHTTPContext) error {
	schemaTokenSvc := services.NewSchemaTokenService(c)
	schemaToken, ierr := schemaTokenSvc.Find(c.Param("repository_id"), c.Param("token_id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, schemaToken)
}
func (sc SchemaController) TokenCreate(c core.IHTTPContext) error {

	input := &requests.SchemaTokenCreate{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}
	schemaTokenSvc := services.NewSchemaTokenService(c)
	schemaToken, ierr := schemaTokenSvc.Create(c.Param("repository_id"), &services.CreateSchemaTokenPayload{
		Name: utils.GetString(input.Name),
		Role: utils.GetString(input.Role),
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, schemaToken)
}
func (sc SchemaController) TokenUpdate(c core.IHTTPContext) error {

	input := &requests.SchemaTokenUpdate{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}
	schemaTokenSvc := services.NewSchemaTokenService(c)
	schemaToken, ierr := schemaTokenSvc.Update(c.Param("repository_id"), &services.UpdateSchemaTokenPayload{
		ID:   c.Param("token_id"),
		Name: utils.GetString(input.Name),
		Role: utils.GetString(input.Role),
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, schemaToken)
}
func (sc SchemaController) TokenDelete(c core.IHTTPContext) error {

	schemaTokenSvc := services.NewSchemaTokenService(c)
	ierr := schemaTokenSvc.Delete(c.Param("repository_id"), c.Param("token_id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.NoContent(http.StatusNoContent)
}
