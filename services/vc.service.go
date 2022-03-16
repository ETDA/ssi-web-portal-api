package services

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/emsgs"
	"gitlab.finema.co/finema/etda/web-portal-api/helpers"
	"gitlab.finema.co/finema/etda/web-portal-api/models"
	"gitlab.finema.co/finema/etda/web-portal-api/views"
	"gorm.io/gorm"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/errmsgs"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type VCVerifyVCPayload struct {
	Message string `json:"message"`
}

type VCVerifyVPPayload struct {
	Message string `json:"message"`
}

type VCRevokePayload struct {
	CID string `json:"cid"`
}

type VCAddPayload struct {
	CID string `json:"cid"`
}

type VCQRVerifyPayload struct {
	CIDs []string
}

type VCJWTMessageHeader struct {
	Alg string `json:"alg"`
	Kid string `json:"kid"`
	Typ string `json:"typ"`
}

type VCJWTMessageClaims struct {
	Exp   int64                 `json:"exp,omitempty"`
	Iat   int64                 `json:"iat,omitempty"`
	Iss   string                `json:"iss,omitempty"`
	Jti   string                `json:"jti,omitempty"`
	Nbf   int64                 `json:"nbf,omitempty"`
	Nonce string                `json:"nonce,omitempty"`
	Sub   string                `json:"sub,omitempty"`
	Aud   string                `json:"aud,omitempty"`
	VC    *VCJWTMessageClaimsVC `json:"vc,omitempty"`
	VP    *VCJWTMessageClaimsVP `json:"vp,omitempty"`
}

type VCJWTMessageClaimsVC struct {
	Context           []string                      `json:"@context"`
	Type              []string                      `json:"type"`
	CredentialSubject core.Map                      `json:"credentialSubject"`
	CredentialSchema  *VCJWTMessageCredentialSchema `json:"credentialSchema"`
}

type VCJWTMessageCredentialSchema struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

type VCJWTMessageClaimsVP struct {
	Context              []string `json:"@context"`
	Type                 []string `json:"type"`
	VerifiableCredential []string `json:"verifiableCredential"`
}

type VCSignRequestPayload struct {
	SchemaName           string
	Signer               string
	Holder               string
	CredentialSubject    core.Map
	CredentialSchemaID   string
	CredentialSchemaType string
	CreatorID            string
}
type VCApprovePayload struct {
	JWT string `json:"jwt"`
}

type VCRejectPayload struct {
	RejectedReason string `json:"rejected_reason"`
}

type VCValidateSchemaPayload struct {
	SchemaID string          `json:"schema_id"`
	Document json.RawMessage `json:"document"`
}

type IVCService interface {
	CreateSignRequest(payload *VCSignRequestPayload) (*models.VC, core.IError)
	Find(id string) (*models.VC, core.IError)
	FindRawVC(id string) (*models.VC, core.IError)
	UpdateSignRequest(status string, jwt string, id string) (*models.VC, core.IError)
	Pagination(status string, startDate string, endDate string, options *core.PageOptions) ([]models.VC, *core.PageResponse, core.IError)
	PaginationByDID(did string, pageOptions *core.PageOptions) ([]models.VC, *core.PageResponse, core.IError)
	FindByCID(cid string) (*models.VC, core.IError)
	VerifyVC(payload *VCVerifyVCPayload) (*views.VCVerify, core.IError)
	VerifyVP(payload *VCVerifyVPPayload) (*views.VPVerify, core.IError)
	GetVCStatus(cid string) (*views.VCStatus, core.IError)
	CreateQRToken(cids []string, didAddress string) (*models.VCQRToken, core.IError)
	FindQRToken(tokenID string) (*models.VCQRToken, core.IError)
	DeleteQRToken(tokenID string) core.IError
	RegisterByServer(orgID string) (*views.VC, core.IError)
	RevokeByServer(orgID string, cid string) core.IError
	Update(payload *VCUpdatePayload) (*models.VC, core.IError)
	Approve(id string, payload *VCApprovePayload) (*models.VC, core.IError)
	Reject(id string, payload *VCRejectPayload) (*models.VC, core.IError)
	VCSigning()
}

type vcService struct {
	ctx               core.IContext
	didSvc            IDIDService
	keySvc            IKeyService
	organizationSvc   IOrganizationService
	mobileUserService IMobileUserService
}

func NewVCService(ctx core.IContext) IVCService {
	return &vcService{
		ctx:               ctx,
		didSvc:            NewDIDService(ctx),
		keySvc:            NewKeyService(ctx),
		organizationSvc:   NewOrganizationService(ctx),
		mobileUserService: NewMobileUserService(ctx),
	}
}

func (s vcService) FindByCID(cid string) (*models.VC, core.IError) {
	vc := &models.VC{}
	err := s.ctx.DB().Preload("Creator").First(vc, "cid = ?", cid).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(errmsgs.NotFound, errmsgs.NotFound)
	}
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	org, ierr := s.organizationSvc.First()
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	if vc.Issuer == utils.GetString(org.DIDAddress) {
		vc.Issuer = org.Name
	}
	user, ierr := s.mobileUserService.FindByDID(vc.Issuer)
	if ierr == nil {
		vc.Issuer = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	}

	user, ierr = s.mobileUserService.FindByDID(vc.Holder)
	if ierr == nil {
		vc.Holder = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	}
	return vc, nil
}

func (s vcService) CreateSignRequest(payload *VCSignRequestPayload) (*models.VC, core.IError) {
	if payload.Signer == "" {
		return s.signByServer(payload)
	}
	return s.signByClient(payload)
}

func (s vcService) findNonceLock(id string) (*models.NonceLock, core.IError) {
	nonceLock := &models.NonceLock{}
	err := s.ctx.DB().First(nonceLock, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(err, errmsgs.NotFound)
	}

	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return nonceLock, nil
}
func (s vcService) FindFirstRequest() (*models.NonceLock, core.IError) {
	nonceLock := &models.NonceLock{}
	err := s.ctx.DB().Model(models.NonceLock{}).Where("is_done = ?", false).Order("created_at asc").First(nonceLock).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	return nonceLock, nil
}
func (s vcService) checkQueue(id string) core.IError {
	nonceLock, ierr := s.findNonceLock(id)
	if ierr != nil {
		if !errmsgs.IsNotFoundError(ierr) {
			return s.ctx.NewError(ierr, ierr)
		}
	}
	var count int64
	s.ctx.DB().Model(models.NonceLock{}).Where("created_at <=? AND is_done = ? AND id != ?", nonceLock.CreatedAt, false, nonceLock.ID).Count(&count)
	if count > 0 {
		return s.ctx.NewError(emsgs.VCSigningUnavailable, emsgs.VCSigningUnavailable)
	}
	return nil
}

func (s vcService) unlockQueue(id string) core.IError {
	nonceLock, ierr := s.findNonceLock(id)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}
	nonceLock.IsDone = true
	nonceLock.UpdatedAt = utils.GetCurrentDateTime()
	err := s.ctx.DB().Updates(nonceLock).Error
	if err != nil {
		return s.ctx.NewError(err, errmsgs.DBError)
	}
	return nil
}

func (s vcService) waitAndLockNonce(id string, interval time.Duration, maxRetry int) core.IError {
	nonceLock, ierr := s.findNonceLock(id)
	if ierr != nil {
		if !errmsgs.IsNotFoundError(ierr) {
			return s.ctx.NewError(ierr, ierr)
		}
	}

	retired := 0
	for retired <= maxRetry {
		err := s.checkQueue(nonceLock.ID)
		if err == nil {
			return nil
		}
		time.Sleep(interval)
		retired++

	}

	// cannot queue in after time interval
	return s.ctx.NewError(emsgs.NonceIsLocked, emsgs.NonceIsLocked)
}
func (s vcService) createNonceLock(vcID string) (*models.NonceLock, core.IError) {
	nonceLock := models.NewNonceLock(vcID)
	err := s.ctx.DB().Create(nonceLock).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	return nonceLock, nil
}

func (s vcService) validateSchema(payload *VCValidateSchemaPayload) (*views.ValidateSchema, core.IError) {
	result := &views.ValidateSchema{}

	ierr := core.RequesterToStruct(result, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Post("/schemas/validate", payload, &core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVVCSchemaServiceBaseURL),
		})
	})
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	return result, nil
}

func (s vcService) signByServer(payload *VCSignRequestPayload) (*models.VC, core.IError) {
	org, ierr := s.organizationSvc.First()
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	holder, ierr := s.didSvc.Find(payload.Holder)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	orgDID := utils.GetString(org.DIDAddress)
	signer, ierr := s.didSvc.Find(orgDID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	nonce, ierr := s.didSvc.GetNonce(signer.ID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	result, ierr := s.validateSchema(&VCValidateSchemaPayload{
		Document: json.RawMessage(utils.StructToString(payload.CredentialSubject)),
		SchemaID: payload.CredentialSchemaID,
	})
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	if !result.Valid {
		return nil, s.ctx.NewError(emsgs.JWTInValid, emsgs.JWTInValid)
	}
	now := utils.GetCurrentDateTime()
	claims := &VCJWTMessageClaims{
		Iss:   signer.ID,
		Nonce: nonce,
		Sub:   holder.ID,
		VC: &VCJWTMessageClaimsVC{
			Context:           []string{consts.VCContext},
			Type:              []string{consts.VCType, payload.CredentialSchemaType},
			CredentialSubject: payload.CredentialSubject,
			CredentialSchema: &VCJWTMessageCredentialSchema{
				ID:   payload.CredentialSchemaID,
				Type: consts.SchemaCredentialSchemaType,
			},
		},
	}

	claimsBase64 := helpers.JSONToBase64NoPadding(claims)
	vcModels := &models.VC{
		ID:           utils.GetUUID(),
		SchemaName:   payload.SchemaName,
		SchemaType:   payload.CredentialSchemaType,
		IssuanceDate: now,
		Issuer:       claims.Iss,
		Holder:       claims.Sub,
		CreatorID:    payload.CreatorID,
		JWT:          claimsBase64,
		Status:       consts.VCStatusPending,
	}

	err := s.ctx.DB().Create(vcModels).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	_, ierr = s.createNonceLock(vcModels.ID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	return s.Find(vcModels.ID)
}
func (s vcService) vcSign(id string) core.IError {
	vcModels, ierr := s.FindRawVC(id)

	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}
	claimJSON, _ := base64.RawStdEncoding.DecodeString(vcModels.JWT)
	claimVC := VCJWTMessageClaims{}
	utils.JSONParse(claimJSON, &claimVC)
	org, ierr := s.organizationSvc.First()
	vc, ierr := s.RegisterByServer(org.ID)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}

	_, ierr = s.didSvc.Find(vcModels.Holder)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}

	nonce, ierr := s.didSvc.GetNonce(vcModels.Issuer)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}

	orgDIDDocument, ierr := s.didSvc.Find(vcModels.Issuer)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}
	header := &VCJWTMessageHeader{
		Alg: "ES256",
		Typ: "JWT",
		Kid: orgDIDDocument.VerificationMethod[0].ID,
	}

	now := utils.GetCurrentDateTime()
	claims := &VCJWTMessageClaims{
		Iss:   vcModels.Issuer,
		Jti:   vc.CID,
		Nbf:   now.Unix(),
		Nonce: nonce,
		Sub:   vcModels.Holder,
		VC:    claimVC.VC,
	}

	headerBase64 := helpers.JSONToBase64URLNoPadding(header)
	claimsBase64 := helpers.JSONToBase64URLNoPadding(claims)

	message := fmt.Sprintf("%s.%s", headerBase64, claimsBase64)

	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}
	signature, ierr := s.keySvc.Sign(utils.GetString(org.EncryptedID), message)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}
	jwt := fmt.Sprintf("%s.%s", message, signature.Signature)
	ierr = s.activeByServer(org.ID, vc.CID, utils.NewSha256(jwt))
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}
	vcModels.CID = claims.Jti
	vcModels.IssuanceDate = now
	vcModels.JWT = jwt
	vcModels.Status = consts.VCStatusActive
	err := s.ctx.DB().Updates(vcModels).Error
	if err != nil {
		return s.ctx.NewError(err, errmsgs.DBError)
	}
	return nil
}

// func (s vcService) signByServer(payload *VCSignRequestPayload) (*models.VC, core.IError) {
// 	org, ierr := s.organizationSvc.First()
// 	if ierr != nil {
// 		return nil, s.ctx.NewError(ierr, ierr)
// 	}
// 	nonceLock, ierr := s.createNonceLock()
// 	if ierr != nil {
// 		return nil, s.ctx.NewError(ierr, ierr)
// 	}
// 	orgDID := utils.GetString(org.DIDAddress)
// 	sessionID := nonceLock.ID
// 	ierr = s.waitAndLockNonce(sessionID, 500*time.Millisecond, 50)
// 	if ierr != nil {
// 		s.unlockQueue(sessionID)
// 		return nil, s.ctx.NewError(emsgs.VCSigningUnavailable, emsgs.VCSigningUnavailable)
// 	}
// 	vc, ierr := s.RegisterByServer(org.ID)
// 	if ierr != nil {
// 		nonceError := s.unlockQueue(sessionID)
// 		if nonceError != nil {
// 			return nil, s.ctx.NewError(nonceError, nonceError)
// 		}
//
// 		return nil, s.ctx.NewError(ierr, ierr)
// 	}
//
// 	holder, ierr := s.didSvc.Find(payload.Holder)
// 	if ierr != nil {
// 		nonceError := s.unlockQueue(sessionID)
// 		if nonceError != nil {
// 			return nil, s.ctx.NewError(nonceError, nonceError)
// 		}
//
// 		return nil, s.ctx.NewError(ierr, ierr)
// 	}
//
// 	nonce, ierr := s.didSvc.GetNonce(orgDID)
// 	if ierr != nil {
// 		nonceError := s.unlockQueue(sessionID)
// 		if nonceError != nil {
// 			return nil, s.ctx.NewError(nonceError, nonceError)
// 		}
//
// 		return nil, s.ctx.NewError(ierr, ierr)
// 	}
//
// 	orgDIDDocument, ierr := s.didSvc.Find(orgDID)
// 	if ierr != nil {
// 		nonceError := s.unlockQueue(sessionID)
// 		if nonceError != nil {
// 			return nil, s.ctx.NewError(nonceError, nonceError)
// 		}
//
// 		return nil, s.ctx.NewError(ierr, ierr)
// 	}
// 	header := &VCJWTMessageHeader{
// 		Alg: "ES256",
// 		Typ: "JWT",
// 		Kid: orgDIDDocument.VerificationMethod[0].ID,
// 	}
//
// 	now := utils.GetCurrentDateTime()
// 	claims := &VCJWTMessageClaims{
// 		Iss:   orgDID,
// 		Jti:   vc.CID,
// 		Nbf:   now.Unix(),
// 		Nonce: nonce,
// 		Sub:   holder.ID,
// 		VC: &VCJWTMessageClaimsVC{
// 			Context:           []string{consts.VCContext},
// 			Type:              []string{consts.VCType, payload.CredentialSchemaType},
// 			CredentialSubject: payload.CredentialSubject,
// 			CredentialSchema: &VCJWTMessageCredentialSchema{
// 				ID:   payload.CredentialSchemaID,
// 				Type: consts.SchemaCredentialSchemaType,
// 			},
// 		},
// 	}
//
// 	headerBase64 := helpers.JSONToBase64URLNoPadding(header)
// 	claimsBase64 := helpers.JSONToBase64URLNoPadding(claims)
//
// 	message := fmt.Sprintf("%s.%s", headerBase64, claimsBase64)
//
// 	signature, ierr := s.keySvc.Sign(utils.GetString(org.EncryptedID), message)
// 	if ierr != nil {
// 		nonceError := s.unlockQueue(sessionID)
// 		if nonceError != nil {
// 			return nil, s.ctx.NewError(nonceError, nonceError)
// 		}
//
// 		return nil, s.ctx.NewError(ierr, ierr)
// 	}
// 	jwt := fmt.Sprintf("%s.%s", message, signature.Signature)
// 	ierr = s.activeByServer(org.ID, vc.CID, utils.NewSha256(jwt))
// 	if ierr != nil {
// 		nonceError := s.unlockQueue(sessionID)
// 		if nonceError != nil {
// 			return nil, s.ctx.NewError(nonceError, nonceError)
// 		}
//
// 		return nil, s.ctx.NewError(ierr, ierr)
// 	}
//
// 	nonceError := s.unlockQueue(sessionID)
// 	if nonceError != nil {
// 		return nil, s.ctx.NewError(nonceError, nonceError)
// 	}
//
// 	result, ierr := s.VerifyVC(&VCVerifyVCPayload{
// 		Message: jwt,
// 	})
//
// 	if ierr != nil {
// 		return nil, s.ctx.NewError(ierr, ierr)
// 	}
// 	if result.VerificationResult != true {
// 		return nil, s.ctx.NewError(emsgs.JWTInValid, emsgs.JWTInValid)
// 	}
// 	vcModels := &models.VC{
// 		ID:           utils.GetUUID(),
// 		CID:          vc.CID,
// 		SchemaName:   payload.SchemaName,
// 		SchemaType:   payload.CredentialSchemaType,
// 		IssuanceDate: now,
// 		Issuer:       claims.Iss,
// 		Holder:       claims.Sub,
// 		CreatorID:    payload.CreatorID,
// 		JWT:          jwt,
// 		Status:       consts.VCStatusActive,
// 	}
//
// 	err := s.ctx.DB().Create(vcModels).Error
// 	if err != nil {
// 		return nil, s.ctx.NewError(err, errmsgs.DBError)
// 	}
//
// 	return s.FindByCID(vcModels.CID)
// }

func (s vcService) signByClient(payload *VCSignRequestPayload) (*models.VC, core.IError) {

	holder, ierr := s.didSvc.Find(payload.Holder)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	signer, ierr := s.didSvc.Find(payload.Signer)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	nonce, ierr := s.didSvc.GetNonce(signer.ID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	now := utils.GetCurrentDateTime()
	claims := &VCJWTMessageClaims{
		Iss:   payload.Signer,
		Nonce: nonce,
		Sub:   holder.ID,
		VC: &VCJWTMessageClaimsVC{
			Context:           []string{consts.VCContext},
			Type:              []string{consts.VCType, payload.CredentialSchemaType},
			CredentialSubject: payload.CredentialSubject,
			CredentialSchema: &VCJWTMessageCredentialSchema{
				ID:   payload.CredentialSchemaID,
				Type: consts.SchemaCredentialSchemaType,
			},
		},
	}

	claimsBase64 := helpers.JSONToBase64NoPadding(claims)

	vcModels := &models.VC{
		ID:           utils.GetUUID(),
		SchemaName:   payload.SchemaName,
		SchemaType:   payload.CredentialSchemaType,
		IssuanceDate: now,
		Issuer:       claims.Iss,
		Holder:       claims.Sub,
		CreatorID:    payload.CreatorID,
		JWT:          claimsBase64,
		Status:       consts.VCStatusPending,
	}

	org, ierr := s.organizationSvc.First()
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	notificationItem := NotificationItem{
		Title:       "New VC signing request",
		Body:        "มีคำขอออกเอกสารมาถึงคุณ",
		ImageURL:    "https://ssi-test.teda.th/logo.png",
		Category:    "",
		Icon:        "",
		ClickAction: "",
		Sound:       "",
		Priority:    "high",
		Data: map[string]string{
			"title":            "New VC Signing Rquest",
			"body":             "มีคำขอออกเอกสารมาถึงคุณ",
			"message":          vcModels.JWT,
			"creator":          org.Name,
			"approve_endpoint": fmt.Sprintf("%s/api/web/vcs/%s/approve", s.ctx.ENV().String(consts.ENVWebPortalBaseURL), vcModels.ID),
			"reject_endpoint":  fmt.Sprintf("%s/api/web/vcs/%s/reject", s.ctx.ENV().String(consts.ENVWebPortalBaseURL), vcModels.ID),
		},
	}
	ierr = s.mobileUserService.SendNotification(&SendNotificationPayload{
		DidAddress:    payload.Signer,
		Notifications: []NotificationItem{notificationItem},
	})

	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	err := s.ctx.DB().Create(vcModels).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return s.Find(vcModels.ID)
}

func (s vcService) Find(id string) (*models.VC, core.IError) {
	vc := &models.VC{}
	err := s.ctx.DB().Preload("Creator").First(vc, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(errmsgs.NotFound, errmsgs.NotFound)
	}
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	org, ierr := s.organizationSvc.First()
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	if vc.Issuer == utils.GetString(org.DIDAddress) {
		vc.Issuer = org.Name
	}
	user, ierr := s.mobileUserService.FindByDID(vc.Issuer)
	if ierr == nil {
		vc.Issuer = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	}

	user, ierr = s.mobileUserService.FindByDID(vc.Holder)
	if ierr == nil {
		vc.Holder = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
	}
	vcStatus, ierr := s.GetVCStatus(vc.CID)
	if ierr == nil {
		if vcStatus.Status != nil {
			if utils.GetString(vcStatus.Status) == consts.VCStatusBlockchainRevoke {
				vc.Status = consts.VCStatusRevoked
			}
		}
	}

	return vc, nil
}

func (s vcService) FindRawVC(id string) (*models.VC, core.IError) {
	vc := &models.VC{}
	err := s.ctx.DB().Preload("Creator").First(vc, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(errmsgs.NotFound, errmsgs.NotFound)
	}
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	return vc, nil
}

func (s vcService) UpdateSignRequest(status string, jwt string, id string) (*models.VC, core.IError) {
	vc, ierr := s.Find(id)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	if vc.Status == consts.VCStatusCanceled {
		return nil, s.ctx.NewError(emsgs.VCStatusCanceledError, emsgs.VCStatusCanceledError)
	}

	if jwt != "" && len(strings.Split(jwt, ".")) != 3 {
		return nil, s.ctx.NewError(emsgs.VCJWTInvalidFormError, emsgs.VCJWTInvalidFormError)
	}

	vc.Status = status
	if jwt != "" {
		vc.JWT = jwt
	}

	err := s.ctx.DB().Updates(vc).Error
	if err != nil {
		return nil, s.ctx.NewError(errmsgs.DBError, errmsgs.DBError)
	}

	return s.Find(vc.ID)
}

func (s vcService) Pagination(status string, startDate string, endDate string, pageOptions *core.PageOptions) ([]models.VC, *core.PageResponse, core.IError) {
	items := make([]models.VC, 0)
	db := s.ctx.DB().Preload("Creator")
	var t time.Time
	const layout = "2006-01-02"
	if len(pageOptions.OrderBy) == 0 {
		pageOptions.OrderBy = []string{"issuance_date desc"}
	}

	if status != "" && status != "ALL" {
		db = db.Where("status = ?", status)
	}

	if startDate != "" {
		t, _ = time.Parse(layout, startDate)
		db = db.Where("issuance_date >= ?", t)
	}

	if endDate != "" {
		t, _ = time.Parse(layout, endDate)
		t = t.AddDate(0, 0, 1)
		db = db.Where("issuance_date <= ?", t)
	}

	db = core.SetSearchSimple(db, pageOptions.Q, []string{"schema_name", "schema_type"})
	pageRes, err := core.Paginate(db, &items, pageOptions)

	if err != nil {
		return nil, nil, s.ctx.NewError(errmsgs.DBError, errmsgs.DBError)
	}
	org, ierr := s.organizationSvc.First()
	if ierr != nil {
		return nil, nil, s.ctx.NewError(ierr, ierr)
	}
	for index, item := range items {
		if item.Issuer == utils.GetString(org.DIDAddress) {
			items[index].Issuer = org.Name
		}
		user, ierr := s.mobileUserService.FindByDID(item.Issuer)
		if ierr == nil {
			items[index].Issuer = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
		}

		user, ierr = s.mobileUserService.FindByDID(item.Holder)
		if ierr == nil {
			items[index].Holder = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
		}
		vcStatus, ierr := s.GetVCStatus(item.CID)
		if ierr == nil {
			if vcStatus.Status != nil {
				if utils.GetString(vcStatus.Status) == consts.VCStatusBlockchainRevoke {
					items[index].Status = consts.VCStatusRevoked
				}
			}
		}
	}
	return items, pageRes, nil
}

func (s vcService) PaginationByDID(did string, pageOptions *core.PageOptions) ([]models.VC, *core.PageResponse, core.IError) {
	items := make([]models.VC, 0)
	db := s.ctx.DB().Preload("Creator")
	if len(pageOptions.OrderBy) == 0 {

		pageOptions.OrderBy = []string{"issuance_date desc"}

	}
	if did != "" {
		db = db.Where("holder = ?", did)
	}
	db = db.Where("status = ?", consts.VCStatusActive)
	pageRes, err := core.Paginate(db, &items, pageOptions)

	if err != nil {
		return nil, nil, s.ctx.NewError(errmsgs.DBError, errmsgs.DBError)
	}

	org, ierr := s.organizationSvc.First()
	if ierr != nil {
		return nil, nil, s.ctx.NewError(ierr, ierr)
	}

	for index, item := range items {
		if item.Issuer == utils.GetString(org.DIDAddress) {
			items[index].Issuer = org.Name
		}
		user, ierr := s.mobileUserService.FindByDID(item.Issuer)
		if ierr == nil {
			items[index].Issuer = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
		}

		user, ierr = s.mobileUserService.FindByDID(item.Holder)
		if ierr == nil {
			items[index].Holder = fmt.Sprintf("%s %s", user.FirstName, user.LastName)
		}
	}
	return items, pageRes, nil
}

func (s vcService) VerifyVC(payload *VCVerifyVCPayload) (*views.VCVerify, core.IError) {
	view := &views.VCVerify{}
	ierr := core.RequesterToStruct(view, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Post("/vc/verify", core.Map{
			"message": payload.Message,
		}, &core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVVerifyServiceBaseURL),
		})
	})

	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	return view, nil
}

func (s vcService) VerifyVP(payload *VCVerifyVPPayload) (*views.VPVerify, core.IError) {
	view := &views.VPVerify{}
	ierr := core.RequesterToStruct(view, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Post("/vp/verify", core.Map{
			"message": payload.Message,
		}, &core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVVerifyServiceBaseURL),
		})
	})

	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	return view, nil
}
func (s vcService) GetVCStatus(cid string) (*views.VCStatus, core.IError) {
	view := &views.VCStatus{}
	ierr := core.RequesterToStruct(view, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Get(fmt.Sprintf("/vc/status/%s", cid), &core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVVCStatusServiceBaseURL),
		})
	})

	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	return view, nil
}

func (s vcService) CreateQRToken(cids []string, didAddress string) (*models.VCQRToken, core.IError) {
	qrToken := models.NewQRToken(cids)
	qrToken.DIDAddress = didAddress
	err := s.ctx.DB().Create(qrToken).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return qrToken, nil
}

func (s vcService) FindQRToken(tokenID string) (*models.VCQRToken, core.IError) {
	qrToken := &models.VCQRToken{}
	err := s.ctx.DB().First(qrToken, "id = ?", tokenID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, s.ctx.NewError(errmsgs.NotFound, errmsgs.NotFound)
	}
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	return qrToken, nil
}

func (s vcService) DeleteQRToken(tokenID string) core.IError {
	qrToken, ierr := s.FindQRToken(tokenID)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}
	err := s.ctx.DB().Delete(qrToken).Error
	if ierr != nil {
		return s.ctx.NewError(err, errmsgs.DBError)
	}
	return nil
}

func (s vcService) RegisterByServer(orgID string) (*views.VC, core.IError) {
	org, ierr := s.organizationSvc.Find(orgID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	nonce, ierr := s.didSvc.GetNonce(utils.GetString(org.DIDAddress))
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	message := core.Map{
		"operation":   consts.OperationVCRegister,
		"did_address": utils.GetString(org.DIDAddress),
		"nonce":       nonce,
	}

	keySign, ierr := s.keySvc.SignJSON(utils.GetString(org.EncryptedID), message)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	view := &views.VC{}
	ierr = core.RequesterToStruct(view, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Post("/vc", core.Map{
			"message": keySign.Message,
		}, &core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVDIDServiceBaseURL),
			Headers: map[string][]string{
				"x-signature": {keySign.Signature},
			},
		})
	})

	return view, ierr
}

type VCRevokeResponse struct {
	Result string `json:"result"`
}

func (s vcService) RevokeByServer(orgID string, id string) core.IError {
	vc, ierr := s.Find(id)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}
	org, ierr := s.organizationSvc.Find(orgID)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}

	nonce, ierr := s.didSvc.GetNonce(utils.GetString(org.DIDAddress))
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}

	message := core.Map{
		"status":      consts.VCStatusBlockchainRevoke,
		"cid":         vc.CID,
		"operation":   consts.OperationVCUpdateStatus,
		"did_address": utils.GetString(org.DIDAddress),
		"nonce":       nonce,
	}

	keySign, ierr := s.keySvc.SignJSON(utils.GetString(org.EncryptedID), message)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}

	result := &VCRevokeResponse{}
	ierr = core.RequesterToStruct(result, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Put(fmt.Sprintf("/vc/status/%s", vc.CID), core.Map{
			"message": keySign.Message,
		}, &core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVVCStatusServiceBaseURL),
			Headers: map[string][]string{
				"x-signature": {keySign.Signature},
			},
		})
	})
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}
	vc.Status = consts.VCStatusRevoked
	err := s.ctx.DB().Updates(vc).Error
	if err != nil {
		return s.ctx.NewError(err, errmsgs.DBError)
	}
	return nil
}

func (s vcService) activeByServer(orgID string, cid string, vcHash string) core.IError {

	org, ierr := s.organizationSvc.Find(orgID)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}

	nonce, ierr := s.didSvc.GetNonce(utils.GetString(org.DIDAddress))
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}

	message := core.Map{
		"status":      consts.VCStatusBlockchainActive,
		"cid":         cid,
		"operation":   consts.OperationVCAddStatus,
		"did_address": utils.GetString(org.DIDAddress),
		"vc_hash":     vcHash,
		"nonce":       nonce,
	}

	keySign, ierr := s.keySvc.SignJSON(utils.GetString(org.EncryptedID), message)
	if ierr != nil {
		return s.ctx.NewError(ierr, ierr)
	}

	result := &VCRevokeResponse{}
	ierr = core.RequesterToStruct(result, func() (*core.RequestResponse, error) {
		return s.ctx.Requester().Post("/vc/status", core.Map{
			"message": keySign.Message,
		}, &core.RequesterOptions{
			BaseURL: s.ctx.ENV().String(consts.ENVVCStatusServiceBaseURL),
			Headers: map[string][]string{
				"x-signature": {keySign.Signature},
			},
		})
	})
	return ierr
}
func (s vcService) Update(payload *VCUpdatePayload) (*models.VC, core.IError) {
	vc, ierr := s.Find(payload.ID)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	if payload.Status == consts.VCStatusCanceled && vc.Status != consts.VCStatusPending {
		return nil, s.ctx.NewError(ierr, emsgs.VCStatusUpdateCanceledError)

	}
	vc.Status = payload.Status
	err := s.ctx.DB().Updates(vc).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	return s.Find(vc.ID)
}

type VCUpdatePayload struct {
	ID     string
	Status string
}

func (s vcService) Approve(id string, payload *VCApprovePayload) (*models.VC, core.IError) {
	vc, ierr := s.FindRawVC(id)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	if vc.Status != consts.VCStatusPending {
		return nil, s.ctx.NewError(emsgs.VCSigningStatusExists, emsgs.VCSigningStatusExists)
	}

	result, ierr := s.VerifyVC(&VCVerifyVCPayload{
		Message: payload.JWT,
	})
	if vc.Issuer != result.Issuer {
		return nil, s.ctx.NewError(errmsgs.SignatureInValid, errmsgs.SignatureInValid)
	}
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}
	if result.VerificationResult == false {
		return nil, s.ctx.NewError(emsgs.VCInvalid, emsgs.VCInvalid)
	}

	vc.CID = result.CID
	vc.IssuanceDate = result.IssuanceDate
	vc.JWT = payload.JWT
	vc.Status = consts.VCStatusActive
	err := s.ctx.DB().Updates(vc).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}
	return s.Find(id)
}
func (s vcService) Reject(id string, payload *VCRejectPayload) (*models.VC, core.IError) {
	vc, ierr := s.FindRawVC(id)
	if ierr != nil {
		return nil, s.ctx.NewError(ierr, ierr)
	}

	if vc.Status != consts.VCStatusPending {
		return nil, s.ctx.NewError(emsgs.VCSigningStatusExists, emsgs.VCSigningStatusExists)
	}

	vc.RejectReason = payload.RejectedReason
	vc.Status = consts.VCStatusReject
	err := s.ctx.DB().Updates(vc).Error
	if err != nil {
		return nil, s.ctx.NewError(err, errmsgs.DBError)
	}

	return s.Find(id)
}

func (s vcService) VCSigning() {
	sign := func() bool {
		nonceLock, err := s.FindFirstRequest()
		if err != nil {
			s.ctx.Log().Debug("Don't have the request to sign in Queue")
			return false
		}
		err = s.vcSign(nonceLock.VCID)
		if err != nil {
			vcs, _ := s.FindRawVC(nonceLock.VCID)
			vcs.Status = consts.VCStatusCanceled
			s.ctx.DB().Updates(vcs)
		}
		nonceLock.IsDone = true
		s.ctx.DB().Updates(nonceLock)
		return true

	}
	for {
		if sign() {

			time.Sleep(2 * time.Second)
		} else {
			time.Sleep(10 * time.Second)
		}
	}
}
