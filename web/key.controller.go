package web

import (
	"crypto/ecdsa"
	"crypto/rsa"
	"net/http"
	"strings"

	"gitlab.finema.co/finema/etda/web-portal-api/consts"
	"gitlab.finema.co/finema/etda/web-portal-api/emsgs"
	"gitlab.finema.co/finema/etda/web-portal-api/helpers"
	"gitlab.finema.co/finema/etda/web-portal-api/requests"
	"gitlab.finema.co/finema/etda/web-portal-api/services"
	core "ssi-gitlab.teda.th/ssi/core"
	"ssi-gitlab.teda.th/ssi/core/utils"
)

type KeyController struct{}

func (a KeyController) Generate(c core.IHTTPContext) error {
	orgSvc := services.NewOrganizationService(c)
	org, ierr := orgSvc.First()
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	item, ierr := orgSvc.GenerateKey(org.ID)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusCreated, item)
}

func (a KeyController) GenerateRSA(c core.IHTTPContext) error {
	orgSvc := services.NewOrganizationService(c)
	org, ierr := orgSvc.First()
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}
	item, ierr := orgSvc.GenerateWithRSAKey(org.ID)
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusCreated, item)
}

func (a KeyController) Upload(c core.IHTTPContext) error {
	input := &requests.KeyStore{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	orgSvc := services.NewOrganizationService(c)
	org, ierr := orgSvc.First()
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	item, ierr := orgSvc.StoreKey(org.ID, &services.KeyStorePayload{
		PublicKey:  utils.GetString(input.PublicKey),
		PrivateKey: utils.GetString(input.PrivateKey),
		KeyType:    utils.GetString(input.KeyType),
	})
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusCreated, item)
}

func (a KeyController) UploadX509(c core.IHTTPContext) error {
	input := &requests.KeyUpload{}
	if err := c.BindWithValidate(input); err != nil {
		return c.JSON(err.GetStatus(), err.JSON())
	}

	cert, err := helpers.ReadX509CertificatePEM(utils.GetString(input.X509Certificate))
	if err != nil {
		return c.JSON(emsgs.ParseCertificateError(err).GetStatus(), emsgs.ParseCertificateError(err).JSON())
	}

	orgSvc := services.NewOrganizationService(c)
	org, ierr := orgSvc.First()
	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	if strings.Contains(cert.SignatureAlgorithm.String(), "RSA") {
		key, err := utils.LoadRSAPrivateKey(utils.GetString(input.X509Key))
		if err != nil {
			return c.JSON(emsgs.ParseKeyError(err).GetStatus(), emsgs.ParseKeyError(err).JSON())
		}

		if !cert.PublicKey.(*rsa.PublicKey).Equal(key.Public()) {
			return c.JSON(emsgs.PublicKeyMismatch.GetStatus(), emsgs.PublicKeyMismatch.JSON())
		}

		prvPEM, pblcPEM := utils.EncodeRSAKeyPair(key, &key.PublicKey)

		item, ierr := orgSvc.StoreKey(org.ID, &services.KeyStorePayload{
			PublicKey:  pblcPEM,
			PrivateKey: prvPEM,
			KeyType:    string(consts.KeyTypeRSA),
		})

		if ierr != nil {
			return c.JSON(ierr.GetStatus(), ierr.JSON())
		}

		return c.JSON(http.StatusCreated, item)
	}

	key, err := utils.LoadPrivateKey(utils.GetString(input.X509Key))
	utils.LogStruct(cert.PublicKey.(*ecdsa.PublicKey))
	utils.LogStruct(key.PublicKey)
	if err != nil {
		return c.JSON(emsgs.ParseKeyError(err).GetStatus(), emsgs.ParseKeyError(err).JSON())
	}

	if !cert.PublicKey.(*ecdsa.PublicKey).Equal(key.Public()) {
		return c.JSON(emsgs.PublicKeyMismatch.GetStatus(), emsgs.PublicKeyMismatch.JSON())
	}
	prvPEM, pblcPEM := utils.EncodeKeyPair(key, &key.PublicKey)

	item, ierr := orgSvc.StoreKey(org.ID, &services.KeyStorePayload{
		PublicKey:  pblcPEM,
		PrivateKey: prvPEM,
		KeyType:    string(consts.KeyTypeECDSA),
	})

	if ierr != nil {
		return c.JSON(ierr.GetStatus(), ierr.JSON())
	}

	return c.JSON(http.StatusCreated, item)
}
