package grpcserver

import (
	"context"

	pb "github.com/PythonHacker24/linux-acl-management-aclapi/internal/grpcserver/protos"
)

/* ping handler for gRPC (testing purposes only) */
type PingHandler struct {
	pb.UnimplementedPingServiceServer
}

/* defination for ping handler */
func (h *PingHandler) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{Message: "pong from module"}, nil
}
