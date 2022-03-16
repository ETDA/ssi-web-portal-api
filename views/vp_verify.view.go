package views

import "time"

type VPVerify struct {
	VerificationResult bool       `json:"verification_result"`
	ID                 string     `json:"id"`
	Audience           string     `json:"audience"`
	Issuer             string     `json:"issuer"`
	IssuanceDate       *time.Time `json:"issuance_date"`
	ExpireDate         *time.Time `json:"expiration_date"`
	VC                 []VCVerify `json:"vc"`
}
