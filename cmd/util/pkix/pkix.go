package pkix

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"math/big"
	"time"
)

func CreateSelfSignedCertificate(cn string, orgs []string, dnsNames []string) (key *rsa.PrivateKey, certBytes []byte, err error) {
	certTmpl := &x509.Certificate{
		SerialNumber: big.NewInt(2021),
		Subject: pkix.Name{
			CommonName:   cn,
			Organization: orgs,
		},
		DNSNames:              dnsNames,
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(50, 0, 0),
		IsCA:                  true,
		BasicConstraintsValid: true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageClientAuth,
			x509.ExtKeyUsageServerAuth,
		},
	}

	caKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}
	certBytes, err = x509.CreateCertificate(rand.Reader, certTmpl, certTmpl, &caKey.PublicKey, caKey)
	if err != nil {
		return nil, nil, err
	}

	return caKey, certBytes, nil
}

func CreateDefaultCertificateRequest(cn string, orgs []string, dnsNames []string) (key *rsa.PrivateKey, csr []byte, err error) {
	return CreateCertificateRequest(rand.Reader, 2048, cn, orgs, dnsNames)
}

func CreateCertificateRequest(random io.Reader, bits int, cn string, orgs []string, dnsNames []string) (key *rsa.PrivateKey, csr []byte, err error) {
	key, err = rsa.GenerateKey(random, bits)
	if err != nil {
		return nil, nil, err
	}

	csrTmpl := x509.CertificateRequest{
		Subject: pkix.Name{
			CommonName:   cn,
			Organization: orgs,
		},
		DNSNames:           dnsNames,
		SignatureAlgorithm: x509.SHA256WithRSA,
	}

	csr, err = x509.CreateCertificateRequest(rand.Reader, &csrTmpl, key)
	if err != nil {
		return nil, nil, err
	}

	return key, csr, err
}

func PemPkcs8PKey(privateKey *rsa.PrivateKey) ([]byte, error) {
	pkcs8, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	var pemKey bytes.Buffer
	err = pem.Encode(&pemKey, &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs8,
	})

	if err != nil {
		return nil, err
	}

	return pemKey.Bytes(), nil
}

func PemCertificate(cert []byte) ([]byte, error) {
	return pemCertificate(cert, "CERTIFICATE")
}

func PemCertificateRequest(csr []byte) ([]byte, error) {
	return pemCertificate(csr, "CERTIFICATE REQUEST")
}

func pemCertificate(der []byte, typ string) ([]byte, error) {
	var pemCert bytes.Buffer

	err := pem.Encode(&pemCert, &pem.Block{Type: typ, Bytes: der})
	if err != nil {
		return nil, err
	}

	return pemCert.Bytes(), nil
}
