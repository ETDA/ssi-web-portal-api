package services

import core "ssi-gitlab.teda.th/ssi/core"

type JWTMessage struct {
	core.BaseValidator
	Header    *JWTMessageHeader `json:"Header"`
	Claims    *JWTMessageClaim  `json:"Claims"`
	Signature *string           `json:"Signature"`
}

type JWTMessageHeader struct {
	Alg string `json:"alg"`
	Kid string `json:"kid"`
	Typ string `json:"typ"`
}

type JWTMessageClaim struct {
	Exp   int64               `json:"exp"`
	Iat   int64               `json:"iat"`
	Iss   string              `json:"iss"`
	Jti   string              `json:"jti"`
	Nbf   int64               `json:"nbf"`
	Nonce string              `json:"nonce"`
	Sub   string              `json:"sub"`
	Aud   string              `json:"aud"`
	VC    *JWTMessageClaimsVC `json:"vc"`
	VP    *JWTMessageClaimsVP `json:"vp"`
}

type JWTMessageClaimsVC struct {
	Context           []string                      `json:"@context"`
	Type              []string                      `json:"type"`
	CredentialSubject core.Map                      `json:"credentialSubject"`
	CredentialSchema  *VCJWTMessageCredentialSchema `json:"credentialSchema"`
}

type JWTMessageClaimsVP struct {
	Context              []string `json:"@context"`
	Type                 []string `json:"type"`
	VerifiableCredential []string `json:"verifiableCredential"`
}
