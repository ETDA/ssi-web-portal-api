package emsgs

import (
	core "ssi-gitlab.teda.th/ssi/core"
	"net/http"
)

var KeyAlreadyExists = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "KEY_ALREADY_EXISTS",
	Message: "key already exists",
}

var OrganizationDIDIsEmpty = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "ORGANIZATION_DID_IS_EMPTY",
	Message: "organization's did is empty",
}

var OrganizationPermissionDenied = core.Error{
	Status:  http.StatusForbidden,
	Code:    "ORGANIZATION_PERMISSION_DENIED",
	Message: "organization permission denied",
}
