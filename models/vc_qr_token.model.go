package models

import (
	"encoding/json"
	"gorm.io/gorm"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
	"time"
)

type VCQRToken struct {
	ID         string           `json:"id" gorm:"column:id"`
	CIDs       *json.RawMessage `json:"cids" gorm:"column:cids"`
	Token      string           `json:"token" gorm:"column:token"`
	DIDAddress string           `json:"did_address" gorm:"column:did_address"`
	CreatedAt  *time.Time       `json:"created_at" gorm:"column:created_at"`
	DeletedAt  *gorm.DeletedAt  `json:"deleted_at,omitempty" gorm:"column:deleted_at"`
}

func (m VCQRToken) TableName() string {
	return "vc_qr_tokens"
}

func NewQRToken(cids []string) *VCQRToken {
	id := utils.GetUUID()
	createdAt := utils.GetCurrentDateTime()
	jsonCIDs := json.RawMessage(utils.JSONToString(core.Map{"cids": cids}))
	return &VCQRToken{
		ID:        utils.GetUUID(),
		CIDs:      &jsonCIDs,
		Token:     utils.NewSha256(id + createdAt.String()),
		CreatedAt: createdAt,
	}
}
