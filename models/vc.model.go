package models

import "time"

type VC struct {
	ID           string     `json:"id" gorm:"column:id"`
	CID          string     `json:"cid" gorm:"column:cid"`
	SchemaName   string     `json:"schema_name" gorm:"column:schema_name"`
	SchemaType   string     `json:"schema_type" gorm:"column:schema_type"`
	IssuanceDate *time.Time `json:"issuance_date" gorm:"column:issuance_date"`
	Issuer       string     `json:"issuer" gorm:"column:issuer"`
	Holder       string     `json:"holder" gorm:"column:holder"`
	CreatorID    string     `json:"-" gorm:"creator_id"`
	Creator      *User      `json:"creator,omitempty" gorm:"foreignKey:CreatorID"`
	JWT          string     `json:"jwt" gorm:"column:jwt"`
	Status       string     `json:"status" gorm:"column:status"`
	RejectReason string     `json:"rejected_reason,omitempty" gorm:"column:rejected_reason"`
}

func (m VC) TableName() string {
	return "vcs"
}
