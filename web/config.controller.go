package web

import (
	"net/http"

	"gitlab.finema.co/finema/etda/web-portal-api/requests"
	"gitlab.finema.co/finema/etda/web-portal-api/services"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type ConfigController struct{}

func (cc ConfigController) WalletGet(c core.IHTTPContext) error {
	walletConfigSvc := services.NewWalletConfigService(c)
	configs, ierr := walletConfigSvc.Get()
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, configs)
}

func (cc ConfigController) WalletSetting(c core.IHTTPContext) error {
	input := &requests.WalletConfigCreate{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}
	walletConfigSvc := services.NewWalletConfigService(c)
	config, ierr := walletConfigSvc.Create(&services.WalletConfigCreatePayload{
		Endpoint:    utils.GetString(input.Endpoint),
		AccessToken: utils.GetString(input.AccessToken),
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusCreated, config)
}

func (cc ConfigController) WalletDelete(c core.IHTTPContext) error {
	walletConfigSvc := services.NewWalletConfigService(c)
	ierr := walletConfigSvc.Delete(c.Param("id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.NoContent(http.StatusNoContent)
}
func (cc ConfigController) SchemaRepositoryPagination(c core.IHTTPContext) error {
	service := services.NewSchemaService(c)
	schemaRepositories, pageResponse, ierr := service.PaginationRepository(c.QueryParam("permission"), c.GetPageOptions())
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, core.NewPagination(schemaRepositories, pageResponse))
}

func (cc ConfigController) SchemaRepositoryFind(c core.IHTTPContext) error {
	service := services.NewSchemaService(c)
	schemaRepository, ierr := service.FindRepository(c.Param("id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, schemaRepository)
}

func (cc ConfigController) SchemaRepositoryCreate(c core.IHTTPContext) error {
	input := &requests.VCSchemaRepositoryCreate{}
	ierr := c.BindWithValidate(input)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	service := services.NewSchemaService(c)
	payloads := make([]services.CreateSchemaRepositoryPayload, 0)
	for _, schemaConfig := range input.SchemaConfigs {
		payload := &services.CreateSchemaRepositoryPayload{
			Name:     utils.GetString(schemaConfig.Name),
			Endpoint: utils.GetString(schemaConfig.Endpoint),
			Token:    utils.GetString(schemaConfig.Token),
		}
		payloads = append(payloads, *payload)
	}
	schemaRepository, ierr := service.CreateRepository(payloads)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, schemaRepository)
}

func (cc ConfigController) SchemaRepositoryUpdate(c core.IHTTPContext) error {
	input := &requests.VCSchemaRepositoryUpdate{}
	ierr := c.BindWithValidate(input)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	service := services.NewSchemaService(c)
	schemaRepository, ierr := service.UpdateRepository(&services.UpdateSchemaRepositoryPayload{
		ID:       c.Param("id"),
		Name:     utils.GetString(input.Name),
		Endpoint: utils.GetString(input.Endpoint),
		Token:    utils.GetString(input.Token),
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.JSON(http.StatusOK, schemaRepository)
}

func (cc ConfigController) SchemaRepositoryDelete(c core.IHTTPContext) error {
	service := services.NewSchemaService(c)
	ierr := service.DeleteRepository(c.Param("id"))
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	return c.NoContent(http.StatusNoContent)
}
