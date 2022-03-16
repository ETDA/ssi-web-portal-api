package requests

import (
	"fmt"

	core "ssi-gitlab.teda.th/ssi/core"
)

type VPTagStatus struct {
	core.BaseValidator
	Tags []string `json:"tags"`
}

func (r VPTagStatus) Valid(ctx core.IContext) core.IError {
	if r.Must(r.IsRequiredArray(r.Tags, "tags")) {
		for index, tag := range r.Tags {
			r.Must(r.IsStrRequired(&tag, fmt.Sprintf("tags[%v]", index)))
		}
	}
	return r.Error()
}
