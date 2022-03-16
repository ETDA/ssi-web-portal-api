package views

import (
	"strings"

	"gitlab.finema.co/finema/etda/web-portal-api/models"
)

type VCWithTag struct {
	models.SubmittedVPVC
	Tags string `json:"tags"`
}

func NewVCWithTag(vc *models.SubmittedVPVC, tags []string) *VCWithTag {
	return &VCWithTag{
		SubmittedVPVC: *vc,
		Tags:          strings.Join(tags, ","),
	}
}
