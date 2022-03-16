package views

import "time"

type VCStatus struct {
	CID         string     `json:"cid"`
	DIDAddress  string     `json:"did_address"`
	Status      *string    `json:"status"`
	VCHash      string     `json:"vc_hash"`
	Tags        []string   `json:"tags"`
	ActivatedAt *time.Time `json:"activated_at"`
	RevokedAt   *time.Time `json:"revoked_at"`
}
