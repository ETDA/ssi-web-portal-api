package middlewares

import (
	"net/http"
	"strings"
	"github.com/labstack/echo/v4"
	"gitlab.finema.co/finema/etda/web-portal-api/services"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

func ValidateSignatureMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(core.IHTTPContext)

		if cc.GetSignature() == "" {
			return c.JSON(http.StatusBadRequest, core.NewValidatorFields(core.RequiredM("x-signature")))
		}
		token := strings.TrimSpace(cc.Request().Header.Get("Authorization"))
		vcService := services.NewVCService(cc)

		tokenID := cc.Param("token_id")
		qrToken, ierr := vcService.FindQRToken(tokenID)
		if ierr != nil {
			return c.JSON(ierr.GetStatus(), ierr.JSON())
		}
		didAddress := qrToken.DIDAddress
		didSvc := services.NewDIDService(cc)
		didDocument, ierr := didSvc.Find(didAddress)
		if ierr != nil {
			return c.JSON(ierr.GetStatus(), ierr.JSON())
		}
		isSigValid := false
		for _, key := range didDocument.VerificationMethod {
			isTempSigValid, _ := utils.VerifySignature(
				key.PublicKeyPem,
				cc.GetSignature(),
				token)
			if isTempSigValid {
				isSigValid = true
				break
			}
		}

		if !isSigValid {
			return c.JSON(errmsgs.SignatureInValid.GetStatus(), errmsgs.SignatureInValid.JSON())
		}

		return next(c)
	}
}
