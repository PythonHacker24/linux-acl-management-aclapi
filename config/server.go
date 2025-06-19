package config

import (
	"errors"
	"fmt"

	"github.com/MakeNowJust/heredoc"
)

/* server config for daemon (grpc network configurations) */
type Server struct {
	Host	 	string 	`yaml:"host,omitempty"`
	GrpcPort 	int		`yaml:"grpc_port,omitempty"`	
	TLSEnabled	bool 	`yaml:"tls_enabled,omitempty"`
	TLSCertFile	string 	`yaml:"tls_cert_file,omitempty"`
	TLSKeyFile	string 	`yaml:"tls_key_file,omitempty"`
}

/* normalization function */
func (s *Server) Normalize() error {
	
	/* bind daemon to "0.0.0.0" */
	if s.Host == "" {
		s.Host = "0.0.0.0"
	}

	/* bind gRPC port to 6593 (decided with alphabets closest looking with the word "gRPC") */
	if s.GrpcPort == 0 {
		s.GrpcPort = 6593
	}

	if s.TLSEnabled {
		if s.TLSCertFile == "" {
			return errors.New(heredoc.Doc(`
				TLS certificate file not provided in the config file

				Please check the docs for more information: 
			`))
		}

		if s.TLSKeyFile == "" {
			return errors.New(heredoc.Doc(`
				TLS key file not provided in the config file	

				Please check the docs for more information: 
			`))
		}

	} else {
		/* TLS will be false by default (give a warning) */
		fmt.Printf("Prefer using TLS for security\n\n")
	}

	return nil 
} 
