package web

import (
	"net/http"
	"ssi-gitlab.teda.th/ssi/core"
)

type HomeController struct{}

func (n *HomeController) Get(c core.IHTTPContext) error {
	return c.JSON(http.StatusOK, core.Map{
		"message": "Hello, I'm web portal API",
	})
}
