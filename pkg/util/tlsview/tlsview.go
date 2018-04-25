package tlsview

import (
	"crypto/x509"
	"encoding/pem"
)

func SmallView(tlsSecret []byte) string {
	secret, _ := pem.Decode(tlsSecret)
	cert, err := x509.ParseCertificate(secret.Bytes)
	if err != nil {
		return "set"
	}
	if len(cert.DNSNames) > 0 {
		return cert.DNSNames[0]
	}
	if len(cert.EmailAddresses) > 0 {
		return cert.EmailAddresses[0]
	}
	return ""
}
