package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/services"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type payload struct {
	core.BaseValidator
	Message *string `json:"message"`
}

func (r payload) Valid(ctx core.IContext) core.IError {
	if r.Must(r.IsStrRequired(r.Message, "message")) {
		r.Must(r.IsBase64(r.Message, "message"))
	}

	return r.Error()
}
func ValidateMessageMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(core.IHTTPContext)

		payloadData := &payload{}
		if cc.GetSignature() == "" {
			return c.JSON(http.StatusBadRequest, core.NewValidatorFields(core.RequiredM("x-signature")))
		}
		if err := cc.BindWithValidate(payloadData); err != nil {
			return c.JSON(err.GetStatus(), err.JSON())
		}
		if cc.GetSignature() == "" {
			return c.JSON(http.StatusBadRequest, core.NewValidatorFields(core.RequiredM("x-signature")))
		}
		message := utils.GetString(payloadData.Message)
		_, err := utils.Base64Decode(message) // decode failed
		if err != nil {
			return c.JSON(errmsgs.BadRequest.GetStatus(), errmsgs.BadRequest.JSON())
		}

		vcService := services.NewVCService(cc)
		vc, ierr := vcService.FindRawVC(cc.Param("id"))
		if ierr != nil {
			return c.JSON(ierr.GetStatus(), ierr.JSON())
		}
		didAddress := vc.Issuer
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
				utils.GetString(payloadData.Message))
			if isTempSigValid {
				isSigValid = true
				break
			}
		}

		if !isSigValid {
			return c.JSON(errmsgs.SignatureInValid.GetStatus(), errmsgs.SignatureInValid.JSON())
		}

		c.Set(consts.ContextKeyMessage, message)
		return next(c)
	}
}
