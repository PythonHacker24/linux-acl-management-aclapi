package grpcserver

import (
	"crypto/tls"

	"google.golang.org/grpc/credentials"
)

/* load tls config for gRPC server */
func loadTLSCredentials(TLSCertFile, TLSKeyFile string) (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(TLSCertFile, TLSKeyFile)
	if err != nil {
		return nil, err
	}
	creds := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	})
	return creds, nil
}
