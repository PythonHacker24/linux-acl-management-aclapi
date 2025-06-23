package acl

import (
	"context"
	"encoding/json"
	"net"

	"github.com/PythonHacker24/linux-acl-management-aclapi/config"
	"github.com/PythonHacker24/linux-acl-management-aclapi/internal/grpcserver/protos"
	pb "github.com/PythonHacker24/linux-acl-management-aclapi/internal/grpcserver/protos"
)

/* ACL Server for gRPC endpoint */
type ACLServer struct {
	pb.UnimplementedACLServiceServer
}

/* handler for handling ACL entry requests */
func (s *ACLServer) ApplyACLEntry(ctx context.Context, req *protos.ApplyACLRequest) (*pb.ApplyACLResponse, error) {
	/* set the socket path as per the configuration */
	socketPath := config.APIDConfig.DConfig.SocketPath

	/* create the ACL modification message */
	aclmsg := struct {
		Action string `json:"action"`
		Entry  string `json:"entry"`
		Path   string `json:"path"`
	}{
		Action: req.Entry.Action,
		Entry:  buildACLEntry(req.Entry),
		Path:   req.TargetPath,
	}

	/* marshall the ACL modification message to JSON data */
	acldata, err := json.Marshal(aclmsg)
	if err != nil {
		return &pb.ApplyACLResponse{Success: false, Message: "JSON encoding failed"}, nil
	}

	/* create a unix socket connection to communicate with ACL core daemon */
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return &pb.ApplyACLResponse{Success: false, Message: "Failed to connect to root daemon"}, nil
	}
	defer conn.Close()

	/* write the ACL JSON data into the connection */
	_, err = conn.Write(acldata)
	if err != nil {
		return &pb.ApplyACLResponse{Success: false, Message: "Failed to write to socket"}, nil
	}

	/* max 1KB response from ACL core daemon (CHANGE IF NEEDED) */
	respBuf := make([]byte, 1024)
	aclResp, err := conn.Read(respBuf)
	if err != nil {
		return &pb.ApplyACLResponse{Success: false, Message: "Failed to read from socket"}, nil
	}

	/* create the response for returning back */
	var response struct {
		Success bool   `json:"success"`
		Message string `json:"message"`
	}

	/* unmarshal JSON response */
	err = json.Unmarshal(respBuf[:aclResp], &response)
	if err != nil {
		return &pb.ApplyACLResponse{Success: false, Message: "Failed to parse response"}, nil
	}

	/* send response via gRPC */
	return &pb.ApplyACLResponse{
		Success: response.Success,
		Message: response.Message,
	}, nil
}
