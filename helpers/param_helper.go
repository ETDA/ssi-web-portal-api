package helpers

import (
	core "ssi-gitlab.teda.th/ssi/core"
	xurl "net/url"
	"strconv"
)

func GetParamsFromPageOptions(options *core.PageOptions) xurl.Values {
	params := xurl.Values{}

	if options == nil {
		return params
	}
	if options.Limit != 0 {
		params["limit"] = []string{strconv.FormatInt(options.Limit, 10)}
	}
	if options.Page != 0 {
		params["page"] = []string{strconv.FormatInt(options.Page, 10)}
	}
	if options.Q != "" {
		params["q"] = []string{options.Q}
	}
	if options.OrderBy != nil {
		if len(options.OrderBy) > 0 {
			params["order_by"] = options.OrderBy
		}
	}

	return params
}
