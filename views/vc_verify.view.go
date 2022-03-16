package views

import "time"

type VCVerify struct {
	VerificationResult       bool                   `json:"verification_result"`
	CID                      string                 `json:"cid"`
	Status                   *string                `json:"status"`
	IssuanceDate             *time.Time             `json:"issuance_date"`
	RevokeDate               *time.Time             `json:"revoke_date,omitempty"`
	ExpireDate               *time.Time             `json:"expire_date,omitempty"`
	Type                     []string               `json:"type"`
	Issuer                   string                 `json:"issuer"`
	Holder                   string                 `json:"holder"`
	SchemaVerificationResult map[string]interface{} `json:"schema_validation_result,omitempty"`
	Schema                   interface{}            `json:"schema,omitempty"`
}
