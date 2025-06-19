package aclapi

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
	"crypto/tls"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	pb "github.com/PythonHacker24/linux-acl-management-aclapi/internal/grpcserver/proto"
	"github.com/PythonHacker24/linux-acl-management-aclapi/config"
	"github.com/PythonHacker24/linux-acl-management-aclapi/internal/utils"
)

func main() {
	if err := exec(); err != nil {
		os.Exit(1)	
	}
}

func exec() error {
	
	/* config must load here in exec() if needed to */

	/*
		true for production, false for development mode
		logger is only for gRPC server and core components (after this step)
	*/
	utils.InitLogger(!config.APIDConfig.DConfig.DebugMode)

	/* zap.L() can be used all over the code for global level logging */
	zap.L().Info("Logger Initiated ...")

	/* preparing graceful shutdown */
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-interrupt
		cancel()
	}()

	return run(ctx)
}

func run(ctx context.Context) error {

	address := fmt.Sprintf("%s:%d", 
		config.APIDConfig.Server.Host, 
		config.APIDConfig.Server.GrpcPort,
	)

	var opts []grpc.ServerOption

	if config.APIDConfig.Server.TLSEnabled {
		creds, err := loadTLSCredentials(
			config.APIDConfig.Server.TLSCertFile,
			config.APIDConfig.Server.TLSKeyFile,
		)
		if err != nil {
			zap.L().Error("Failed to load TLS certificate and TLS Key files: %w", err)
			return err
		}

		opts = append(opts, grpc.Creds(creds))
		zap.L().Info("TLS enabled for gRPC")
	} else {
		zap.L().Warn("Proceeding to start gRPC without TLS")
	}

	grpcServer := grpc.NewServer(opts...)
	pb.RegisterACLServiceServer(grpcServer, &ACLServer{})
	
	listener, err := net.Listen("tcp", address)
	if err != nil {
		zap.L().Error("Failed to listen on specified address", 
			zap.String("Address: ", address), 
			zap.Error(err),
		)
		return err
	
	}

	go func() {
		zap.L().Info("gRPC server is listening",
			zap.String("Address", config.APIDConfig.Server.Host),
			zap.Int("Port", config.APIDConfig.Server.GrpcPort),
		)
		if err := grpcServer.Serve(listener); err != nil {
			zap.L().Error("gRPC server error",
				zap.Error(err),
			)
		}
	}()

	<-ctx.Done()

	zap.L().Info("Shutting down gRPC server")

	/* attempting to gracefully shutdown gRPC server */
	gracefulStop := make(chan struct{})
	go func() {
		grpcServer.GracefulStop()
		close(gracefulStop)
	}()

	select {
	case <-gracefulStop:
		zap.L().Info("gRPC server stopped gracefully")
	case <-time.After(5 * time.Second):
		/* gRPC server graceful shutdown timeout reached, need to forcefully shutdown it down */
		zap.L().Info("Graceful shutdown timeout reached. Forcing gRPC server shutdown")
		grpcServer.Stop()
	}

	return nil
}

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
