package helpers

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func ReadX509CertificatePEM(x509PEM string) (*x509.Certificate, error) {
	block, _ := pem.Decode([]byte(x509PEM))
	if block == nil {
		return nil, errors.New("failed to parse certificate PEM")
	}

	return x509.ParseCertificate(block.Bytes)
}
