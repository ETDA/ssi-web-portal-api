package emsgs

import (
	"net/http"

	core "ssi-gitlab.teda.th/ssi/core"
)

var KeyIDIsEmpty = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "KEY_ID_EMPTY",
	Message: "key id must not empty",
}

var KeyMessageIsEmpty = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "MESSAGE_IS_EMPTY",
	Message: "message must not empty",
}

var PublicKeyMismatch = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "PUBLIC_KEY_MISMATCED",
	Message: "the public key in x509 certificate is not matched with private key",
}

var ParseCertificateError = func(err error) core.Error {
	ierr := core.Error{
		Status:  http.StatusBadRequest,
		Code:    "PARSE_CERTIFICATE_FAILED",
		Message: err.Error(),
	}
	return ierr
}

var ParseKeyError = func(err error) core.Error {
	ierr := core.Error{
		Status:  http.StatusBadRequest,
		Code:    "PARSE_KEY_FAILED",
		Message: err.Error(),
	}
	return ierr
}
