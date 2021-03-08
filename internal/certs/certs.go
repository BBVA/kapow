package certs

import (
	"bytes"
	"crypto"
	"crypto/x509"
	"encoding/pem"

	"github.com/BBVA/kapow/internal/logger"
)

type Cert struct {
	X509Cert   *x509.Certificate
	PrivKey    crypto.PrivateKey
	SignedCert []byte
}

func (c Cert) SignedCertPEMBytes() []byte {

	PEM := new(bytes.Buffer)
	err := pem.Encode(PEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: c.SignedCert,
	})
	if err != nil {
		logger.L.Fatal(err)
	}

	return PEM.Bytes()
}
