package pkix

import (
	"crypto/x509"
	"encoding/pem"
	"reflect"
	"testing"
)

func TestCreateSelfSignedCertificate(t *testing.T) {
	var tests = []struct {
		cn       string
		orgs     []string
		dnsNames []string
	}{
		{
			cn:       "local.io",
			dnsNames: nil,
			orgs:     nil,
		},
		{
			cn:       "local.io",
			orgs:     nil,
			dnsNames: []string{"local.io", "*.local.io"},
		},
		{
			cn:       "local.io",
			orgs:     []string{"developers", "Global Security"},
			dnsNames: []string{"local.io", "*.local.io"},
		},
	}
	for _, test := range tests {

		key, cert, err := CreateSelfSignedCertificate(test.cn, test.orgs, test.dnsNames)
		if err != nil {
			t.Error(err)
		}

		xCert, err := x509.ParseCertificate(cert)
		if err != nil {
			t.Error(err)
		}

		err = xCert.CheckSignature(xCert.SignatureAlgorithm, xCert.RawTBSCertificate, xCert.Signature)
		if err != nil {
			t.Error(err)
		}

		if !key.PublicKey.Equal(xCert.PublicKey) {
			t.Error("Public Key not matching: invalid certificate")
		}

		if xCert.Subject.CommonName != test.cn {
			t.Errorf("CommonName: (%q) = %v", test.cn, xCert.Subject.CommonName)
		}

		if !reflect.DeepEqual(xCert.DNSNames, test.dnsNames) {
			t.Errorf("DNSNames: (%q) = %q", test.dnsNames, xCert.DNSNames)
		}

		if !reflect.DeepEqual(xCert.Subject.Organization, test.orgs) {
			t.Errorf("Organization: (%q) = %v", test.orgs, xCert.Subject.Organization)
		}
	}
}

func TestCreateDefaultCertificateRequest(t *testing.T) {
	var tests = []struct {
		cn       string
		orgs     []string
		dnsNames []string
	}{
		{
			cn:       "local.io",
			dnsNames: nil,
			orgs:     nil,
		},
		{
			cn:       "local.io",
			orgs:     nil,
			dnsNames: []string{"local.io", "*.local.io"},
		},
		{
			cn:       "local.io",
			orgs:     []string{"developers", "Global Security"},
			dnsNames: []string{"local.io", "*.local.io"},
		},
	}
	for _, test := range tests {
		key, csr, err := CreateDefaultCertificateRequest(test.cn, test.orgs, test.dnsNames)
		if err != nil {
			t.Error(err)
		}

		xCsr, err := x509.ParseCertificateRequest(csr)
		if err != nil {
			t.Error(err)
		}

		if err = xCsr.CheckSignature(); err != nil {
			t.Errorf("invalid signature: %s", err)
		}

		if xCsr.Subject.CommonName != test.cn {
			t.Errorf("CommonName: (%q) = %v", test.cn, xCsr.Subject.CommonName)
		}

		if !reflect.DeepEqual(xCsr.DNSNames, test.dnsNames) {
			t.Errorf("DNSNames: (%q) = %q", test.dnsNames, xCsr.DNSNames)
		}

		if !reflect.DeepEqual(xCsr.Subject.Organization, test.orgs) {
			t.Errorf("Organization: (%q) = %v", test.orgs, xCsr.Subject.Organization)
		}
		_ = key
	}
}

func TestPemCertificateRequest(t *testing.T) {
	var tests = []struct {
		typ string
		cn  string
	}{
		{
			typ: "CERTIFICATE REQUEST",
			cn:  "local.io",
		},
	}
	for _, test := range tests {
		_, csr, err := CreateDefaultCertificateRequest(test.cn, nil, nil)
		if err != nil {
			t.Error(err)
		}

		pemCsr, err := PemCertificateRequest(csr)
		if err != nil {
			t.Error(err)
		}

		block, _ := pem.Decode(pemCsr)

		if block == nil {
			t.Errorf("pem: codes error")
		}

		if block.Type != test.typ {
			t.Errorf("pem: got %q, wang %v", block.Type, test.typ)
		}

		if reflect.DeepEqual(pemCsr, block.Bytes) {
			t.Error("pem: codes error")
		}
	}
}

func TestPemCertificate(t *testing.T) {
	var tests = []struct {
		typ string
		cn  string
	}{
		{
			typ: "CERTIFICATE",
			cn:  "local.io",
		},
	}
	for _, test := range tests {
		_, csr, err := CreateSelfSignedCertificate(test.cn, nil, nil)
		if err != nil {
			t.Error(err)
		}

		pemCsr, err := PemCertificate(csr)
		if err != nil {
			t.Error(err)
		}

		block, _ := pem.Decode(pemCsr)

		if block == nil {
			t.Errorf("pem: codes error")
		}

		if block.Type != test.typ {
			t.Errorf("pem: got %q, wang %v", block.Type, test.typ)
		}

		if reflect.DeepEqual(pemCsr, block.Bytes) {
			t.Error("pem: codes error")
		}
	}
}
