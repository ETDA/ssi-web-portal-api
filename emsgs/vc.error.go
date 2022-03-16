package emsgs

import (
	"net/http"

	core "ssi-gitlab.teda.th/ssi/core"
)

var VCStatusCanceledError = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "VC_STATUS_CANCELED",
	Message: "VC with canceled status cannot be update",
}

var VCJWTInvalidFormError = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "INVALID_JWT_FORM",
	Message: "VC JWT is invalid form",
}

var VCInvalid = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "INVALID_VC",
	Message: "VC is not valid"}

var VCStatusUpdateCanceledError = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "VC_STATUS_UPDATE",
	Message: "VC's status cannot be update to CANCELED if its status is not PENDING",
}

var VCSigningUnavailable = core.Error{
	Status:  http.StatusUnprocessableEntity,
	Code:    "VC_SIGNING_UNAVAILABLE",
	Message: "Server processing another Signing Request",
}

var VCSigningStatusExists = core.Error{
	Status:  http.StatusUnprocessableEntity,
	Code:    "VC_SIGNING_STATUS_EXISTS",
	Message: "This Document already have status active or reject",
}
