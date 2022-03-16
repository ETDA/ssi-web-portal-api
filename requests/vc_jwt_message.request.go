package requests

import (
	"fmt"
	"strings"

	core "ssi-gitlab.teda.th/ssi/core"
)

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

type VCJWTMessageCredentialSchema struct {
	ID   *string `json:"id"`
	Type *string `json:"type"`
}

type JWTMessageClaimsVP struct {
	Context              []string `json:"@context"`
	Type                 []string `json:"type"`
	VerifiableCredential []string `json:"verifiableCredential"`
}

func (r JWTMessage) Valid(ctx core.IContext) core.IError {
	r.Must(r.IsRequired(r.Header, "Header"))
	r.Must(r.IsRequired(r.Claims, "Claims"))

	r.Must(r.IsStrRequired(&r.Header.Alg, "Header.alg"))
	r.Must(r.IsStrRequired(&r.Header.Typ, "Header.typ"))

	r.Must(r.IsStrRequired(&r.Header.Kid, "Header.kid"))
	r.Must(r.IsStrRequired(&r.Claims.Jti, "Claims.jti"))
	if r.Claims.VC != nil {
		r.Must(r.IsStrRequired(&r.Claims.Sub, "Claims.sub"))
		r.Must(r.IsStrRequired(&r.Claims.Iss, "Claims.iss"))
		r.Must(r.IsRequiredArray(r.Claims.VC.Type, "Claims.vc.type"))
		if r.Must(r.IsArrayMin(r.Claims.VC.Type, 1, "Claims.vc.type")) {
			included := false
			for i, vcType := range r.Claims.VC.Type {
				r.Must(r.IsStrRequired(&vcType, fmt.Sprintf("Claims.vc.type[%v]", i)))
				if vcType == "VerifiableCredential" {
					included = true
				}
			}
			if !included {
				s := strings.Join(r.Claims.VC.Type, "|")
				r.Must(r.IsStrIn(&s, "VerifiableCredential", "Claims.vc.type"))
			}
		}
		r.Must(r.IsRequired(r.Claims.VC.CredentialSubject, "Claims.vc.credentialSubject"))
		if r.Must(r.IsRequired(r.Claims.VC.CredentialSchema, "Claims.vc.credentialSchema")) {
			r.Must(r.IsStrRequired(r.Claims.VC.CredentialSchema.ID, "Claims.vc.credentialSchema.id"))
		}
	}

	if r.Claims.VP != nil {
		r.Must(r.IsStrRequired(&r.Claims.Iss, "Claims.iss"))

		r.Must(r.IsRequiredArray(r.Claims.VP.Type, "Claims.vp.type"))
		if r.Must(r.IsArrayMin(r.Claims.VP.Type, 1, "Claims.vp.type")) {
			included := false
			for i, vpType := range r.Claims.VP.Type {
				r.Must(r.IsStrRequired(&vpType, fmt.Sprintf("Claims.vp.type[%v]", i)))
				if vpType == "VerifiablePresentation" {
					included = true
				}
			}
			if !included {
				s := strings.Join(r.Claims.VP.Type, "|")
				r.Must(r.IsStrIn(&s, "VerifiablePresentation", "Claims.vp.type"))
			}
		}

		r.Must(r.IsStrRequired(&r.Claims.Aud, "Claims.aud"))

		r.Must(r.IsRequiredArray(r.Claims.VP.VerifiableCredential, "Claims.vp.verifiableCredential"))
		r.Must(r.IsArrayMin(r.Claims.VP.VerifiableCredential, 1, "Claims.vp.verifiableCredential"))
	}

	if r.Claims.VC == nil && r.Claims.VP == nil {
		r.Must(false, &core.IValidMessage{
			Name:    "Claims.vc|Claims.vp",
			Code:    "REQUIRED",
			Message: "The Claims.vc or Claims.vp field is required",
		})
	}

	return r.Error()
}
