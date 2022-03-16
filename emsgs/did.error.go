package emsgs

import (
	"net/http"

	core "ssi-gitlab.teda.th/ssi/core"
)

var NonceIsLocked = &core.Error{
	Status:  http.StatusServiceUnavailable,
	Code:    "SERVICE_UNAVAILABLE",
	Message: "nonce with server did is being used at this time. please wait and try again",
}
