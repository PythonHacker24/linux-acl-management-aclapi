package grpcserver

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"

	"google.golang.org/grpc/credentials"
)

/* load tls config for gRPC server */
func loadMutualTLSCredentials(TLSCertFile, TLSKeyFile, TLSCACertFile string) (credentials.TransportCredentials, error) {
	/* load TLS certificate and keys */
	cert, err := tls.LoadX509KeyPair(TLSCertFile, TLSKeyFile)
	if err != nil {
		return nil, err
	}

	/* load CA certificate file */
	caCert, err := os.ReadFile(TLSCACertFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read CA cert: %w", err)
	}
	/* create a new certificate pool and append CA certificate */
	certPool := x509.NewCertPool()
	if ok := certPool.AppendCertsFromPEM(caCert); !ok {
		return nil, fmt.Errorf("failed to add CA cert to pool")
	}

	/* build a TLS config */
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		ClientCAs:    certPool,
		ClientAuth:   tls.RequireAndVerifyClientCert,
		MinVersion:   tls.VersionTLS12,
	}

	/* return TLS credentials */
	return credentials.NewTLS(tlsConfig), nil
}
