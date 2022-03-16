package emsgs

import (
	"net/http"

	core "ssi-gitlab.teda.th/ssi/core"
)

var (
	JWTInValid = core.Error{
		Status:  http.StatusBadRequest,
		Code:    "INVALID_JWT",
		Message: "JWT is not valid"}
)
