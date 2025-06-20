package grpcserver

import (
	"google.golang.org/grpc"

	"github.com/PythonHacker24/linux-acl-management-aclapi/config"
)

/* server struct for gRPC server */
type Server struct {
	GRPC   *grpc.Server
	Config *config.Server
}
