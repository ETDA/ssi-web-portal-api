package emsgs

import (
	"fmt"
	"net/http"

	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
)

var SchemaNotFound = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "SCHEMA_NOT_FOUND",
	Message: "schema is not found",
}

func RequestFormError(err error) core.IError {
	return core.Error{
		Status:  http.StatusBadRequest,
		Code:    "BAD_REQUEST",
		Message: err.Error(),
	}
}

func OpenFileError(err error) core.IError {
	return core.Error{
		Status:  http.StatusInternalServerError,
		Code:    "INTERNAL_SERVER_ERROR",
		Message: err.Error(),
	}
}

func SchemaRepositoryConectError(ierr core.IError) core.IError {
	if ierr.GetCode() == "NETWORK_ERROR" {
		return core.Error{
			Status:  http.StatusBadGateway,
			Code:    "BAD_GATEWAY",
			Message: fmt.Sprintf("requesting url cannot be reached according to the message %s", ierr.OriginalError().Error()),
		}
	}

	if ierr.GetStatus() == http.StatusNotFound {
		return errmsgs.NotFound
	}

	return ierr
}
