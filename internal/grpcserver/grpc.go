package grpcserver

import (
	"fmt"
	"net"

	"github.com/PythonHacker24/linux-acl-management-aclapi/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/PythonHacker24/linux-acl-management-aclapi/internal/acl"
	pb "github.com/PythonHacker24/linux-acl-management-aclapi/internal/grpcserver/protos"
)

func InitServer() (*Server, error) {
	/* gRPC server options */
	var opts []grpc.ServerOption

	/* loading TLS credentials if specified */
	if config.APIDConfig.Server.TLSEnabled {
		creds, err := loadTLSCredentials(
			config.APIDConfig.Server.TLSCertFile,
			config.APIDConfig.Server.TLSKeyFile,
		)
		if err != nil {
			return nil, fmt.Errorf("Failed to load TLS certificate and TLS Key files: %w", err)
		}

		opts = append(opts, grpc.Creds(creds))
		zap.L().Info("TLS enabled for gRPC")
	} else {
		zap.L().Warn("Proceeding to start gRPC without TLS")
	}

	/* setting options to the gRPC server */
	// grpcServer := grpc.NewServer(opts...)
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryServerInterceptor()),
	)

	/* registering services */
	// pb.RegisterACLServiceServer(grpcServer, &ACLServer{})
	pb.RegisterPingServiceServer(grpcServer, &PingHandler{})
	pb.RegisterACLServiceServer(grpcServer, &acl.ACLServer{})

	/* enable reflection if daemon is in debug mode */
	if config.APIDConfig.DConfig.DebugMode {
		reflection.Register(grpcServer)
	}

	return &Server{GRPC: grpcServer, Config: &config.APIDConfig.Server}, nil
}

/* start the gRPC server */
func (s *Server) Start() (net.Listener, error) {

	/* address for gRPC server to bind */
	address := fmt.Sprintf("%s:%d",
		s.Config.Host,
		s.Config.GrpcPort,
	)

	/* creating the gRPC listener */
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return nil, fmt.Errorf("Failed to listen on specified address: %s, error: %w", address, err)
	}

	return listener, err
}
