package server

import (
	"crypto/tls"
	"fmt"
)

func TLSCert(config *tls.Config, certFile, keyFile string) (*tls.Config, error) {
	var err error
	config.Certificates = make([]tls.Certificate, 1)
	config.Certificates[0], err = tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to configure certificates: %w", err)
	}
	return config, nil
}
