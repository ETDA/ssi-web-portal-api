package emsgs

import (
	"net/http"

	core "ssi-gitlab.teda.th/ssi/core"
)

var RequestedVPSubmitError = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "VP_SUBMIT_ERROR",
	Message: "Cannot Submit VP to requested VP with canceled status",
}

var VCNotActiveError = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "VC_INACTIVE_ERROR",
	Message: "Cannot Submit VP with in inactive vc",
}
