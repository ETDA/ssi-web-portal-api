package emsgs

import (
	"net/http"
	core "ssi-gitlab.teda.th/ssi/core"
)

var AuthEmailOrPasswordInvalid = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "INVALID_CREDENTIALS",
	Message: "email or password is an invalid",
}

var AuthCurrentPasswordInvalid = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "INVALID_CURRENT_PASSWORD",
	Message: "current password does not match",
}

var AuthTokenRequired = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "TOKEN_REQUIRED",
	Message: "Authorization header field is required"}

var VCQRTokenRequired = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "TOKEN_REQUIRED",
	Message: "X-Vc-Qr-Token header field is required"}

var AuthTokenInvalid = core.Error{
	Status:  http.StatusUnauthorized,
	Code:    "INVALID_TOKEN",
	Message: "Token is invalid"}

var AuthAPIKeyInvalid = core.Error{
	Status:  http.StatusUnauthorized,
	Code:    "INVALID_API_KEY",
	Message: "Api key is invalid"}
