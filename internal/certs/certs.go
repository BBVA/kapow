package certs

import (
	"bytes"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"

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

func (c Cert) PrivateKeyPEMBytes() []byte {
	PEM := new(bytes.Buffer)
	err := pem.Encode(PEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(c.PrivKey.(*rsa.PrivateKey)),
	})
	if err != nil {
		logger.L.Fatal(err)
	}

	return PEM.Bytes()
}

func GenCert(name, altName string, isServer bool) Cert {

	usage := x509.ExtKeyUsageClientAuth
	if isServer {
		usage = x509.ExtKeyUsageServerAuth
	}

	var dnsNames []string
	var ipAddresses []net.IP
	if altName != "" {
		if ipAddr := net.ParseIP(altName); ipAddr != nil {
			ipAddresses = []net.IP{ipAddr}
		} else {
			dnsNames = []string{altName}
		}
	}

	cert := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		DNSNames:     dnsNames,
		IPAddresses:  ipAddresses,
		Subject: pkix.Name{
			CommonName: name,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  false,
		BasicConstraintsValid: true,
		ExtKeyUsage: []x509.ExtKeyUsage{
			usage,
		},
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		logger.L.Fatal(err)
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, cert, cert, &certPrivKey.PublicKey, certPrivKey)
	if err != nil {
		logger.L.Fatal(err)
	}

	return Cert{
		X509Cert:   cert,
		PrivKey:    certPrivKey,
		SignedCert: certBytes,
	}
}
