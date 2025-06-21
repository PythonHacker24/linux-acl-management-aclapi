package acl

import (
	"fmt"
	pb "github.com/PythonHacker24/linux-acl-management-aclapi/internal/grpcserver/protos"
)

/* builds ACL entry for gRPC handler */
func buildACLEntry(entry *pb.ACLEntry) string {
	prefix := ""
	if entry.IsDefault {
		prefix = "default:"
	}

	entity := entry.Entity
	if entity == "" {
		entity = ""
	}

	return fmt.Sprintf("%s%s:%s:%s", prefix, entry.EntityType, entity, entry.Permissions)
}
