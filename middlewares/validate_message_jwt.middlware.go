package middlewares

import (
	"github.com/labstack/echo/v4"
	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/emsgs"
	"gitlab.finema.co/finema/etda/web-portal-api/helpers"
	"gitlab.finema.co/finema/etda/web-portal-api/requests"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type jwtPayload struct {
	core.BaseValidator
	Message *string `json:"message"`
}

func (r jwtPayload) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsStrRequired(r.Message, "message"))

	return r.Error()
}

func ValidateJWTMessageMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cc := c.(core.IHTTPContext)

		payloadData := &jwtPayload{}
		if err := cc.BindWithValidate(payloadData); err != nil {
			return c.JSON(err.GetStatus(), err.JSON())
		}

		tokenM, _ := helpers.JWTVCDecodingT(utils.GetString(payloadData.Message), []byte(""))
		if tokenM == nil || tokenM.Header == nil || tokenM.Claims == nil || tokenM.Signature == "" {
			return c.JSON(emsgs.JWTInValid.GetStatus(), emsgs.JWTInValid.JSON())
		}

		msgPayload := &requests.JWTMessage{}
		_ = utils.MapToStruct(tokenM, &msgPayload)
		if err := cc.BindWithValidate(msgPayload); err != nil {
			return c.JSON(err.GetStatus(), err.JSON())
		}

		c.Set(consts.ContextKeyJWTData, msgPayload)
		c.Set(consts.ContextKeyJWT, utils.GetString(payloadData.Message))

		return next(c)
	}
}
