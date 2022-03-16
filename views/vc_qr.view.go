package views

import (
	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
)

type VCQRView struct {
	Operation string `json:"operation"`
	Token     string `json:"token"`
	Endpoint  string `json:"endpoint"`
}

func NewVCQRView(vCQRToken *models.VCQRToken, qrCodeURL string) *VCQRView {
	return &VCQRView{
		Endpoint:  qrCodeURL,
		Token:     vCQRToken.Token,
		Operation: consts.OperationGetVC,
	}
}
