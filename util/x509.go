package util

import (
	"crypto/x509"
	"encoding/pem"
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
