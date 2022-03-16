package emsgs

import (
	"net/http"

	core "ssi-gitlab.teda.th/ssi/core"
)

var WalletConfigExists = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "WALLET_CONFIG_EXISTS",
	Message: "Wallet Config is exists",
}
var WalletConfigNotCorrect = core.Error{
	Status:  http.StatusBadRequest,
	Code:    "WALLET_CONFIG_NOT_CORRECT",
	Message: "Endpoint or Access token is not correct",
}
