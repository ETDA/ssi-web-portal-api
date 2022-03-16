package views

import "gitlab.finema.co/finema/etda/web-portal-api/models"

type VC struct {
	CID        string `json:"cid"`
	DIDAddress string `json:"did_address"`
}

type VCItem struct {
	*models.VC
	Tags []string `json:"tags"`
}

func NewVCItem(vc *models.VC, tags []string) *VCItem {
	newTags :=  make([]string, 0)
	if tags != nil {
		newTags = tags
	}
	return &VCItem{
		VC:   vc,
		Tags: newTags,
	}
}
