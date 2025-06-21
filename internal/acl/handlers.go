package acl

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/PythonHacker24/linux-acl-management-aclapi/internal/grpcserver/protos"
	pb "github.com/PythonHacker24/linux-acl-management-aclapi/internal/grpcserver/protos"
)

/* ACL Server for gRPC endpoint */
type ACLServer struct {
	pb.UnimplementedACLServiceServer
}

/* handler for handling ACL entry requests */
func (s *ACLServer) ApplyACLEntry(ctx context.Context, req *protos.ApplyACLRequest) (*pb.ApplyACLResponse, error) {
	entry := req.Entry
	path := req.TargetPath

	/* build the ACL entry */
	aclArg := buildACLEntry(entry)

	/* cmd contains the ACL modification command */
	var cmd *exec.Cmd

	/* execute command as per the request */
	switch entry.Action {
	case "add", "modify":
		cmd = exec.Command("setfacl", "-m", aclArg, path)
	case "remove":
		cmd = exec.Command("setfacl", "-x", aclArg, path)
	default:
		return &pb.ApplyACLResponse{
			Success: false,
			Message: fmt.Sprintf("Unsupported action: %s", entry.Action),
		}, nil
	}

	/* get the output from command execution */
	output, err := cmd.CombinedOutput()
	if err != nil {
		return &pb.ApplyACLResponse{
			Success: false,
			Message: string(output),
		}, nil
	}

	/* return the output */
	return &pb.ApplyACLResponse{
		Success: true,
		Message: string(output),
	}, nil
}
