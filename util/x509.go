package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

func ParseX509Certificates(bytes []byte) ([]*x509.Certificate, error) {
	var certificates []*x509.Certificate

	for len(bytes) > 0 {
		var block *pem.Block
		block, bytes = pem.Decode(bytes)
		if block != nil {
			if certificate, err := x509.ParseCertificate(block.Bytes); err == nil {
				certificates = append(certificates, certificate)
			} else {
				return nil, err
			}
		} else {
			bytes = nil
		}
	}

	return certificates, nil
}

func ParseX509CertPool(bytes []byte) (*x509.CertPool, error) {
	if certificates, err := ParseX509Certificates(bytes); err == nil {
		if len(certificates) > 0 {
			certPool := x509.NewCertPool()
			for _, certificate := range certificates {
				certPool.AddCert(certificate)
			}
			return certPool, nil
		} else {
			return nil, nil
		}
	} else {
		return nil, err
	}
}

var serialNumberLimit = new(big.Int).Lsh(big.NewInt(1), 128)

func CreateSelfSignedX509(organization string, host string) ([]byte, []byte, error) {
	// See: fasthttp

	if serialNumber, err := rand.Int(rand.Reader, serialNumberLimit); err == nil {
		now := time.Now()
		certificate := &x509.Certificate{
			Subject: pkix.Name{
				Organization: []string{organization},
			},
			DNSNames:     []string{host},
			SerialNumber: serialNumber,
			KeyUsage:     x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
			ExtKeyUsage: []x509.ExtKeyUsage{
				x509.ExtKeyUsageServerAuth,
				x509.ExtKeyUsageClientAuth,
			},
			SignatureAlgorithm:    x509.SHA256WithRSA,
			BasicConstraintsValid: true,
			IsCA:                  true,
			NotBefore:             now,
			NotAfter:              now.Add(365 * 24 * time.Hour), // one year
		}

		if key, err := rsa.GenerateKey(rand.Reader, 2048); err == nil {
			if certificateBytes, err := x509.CreateCertificate(rand.Reader, certificate, certificate, &key.PublicKey, key); err == nil {
				return pem.EncodeToMemory(
						&pem.Block{
							Type:  "CERTIFICATE",
							Bytes: certificateBytes,
						},
					), pem.EncodeToMemory(
						&pem.Block{
							Type:  "PRIVATE KEY",
							Bytes: x509.MarshalPKCS1PrivateKey(key),
						},
					), nil
			} else {
				return nil, nil, err
			}
		} else {
			return nil, nil, err
		}
	} else {
		return nil, nil, err
	}
}
