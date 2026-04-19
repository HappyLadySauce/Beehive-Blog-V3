package ctxmeta

import (
	"context"
	"strings"

	"google.golang.org/grpc/metadata"
)

// GetClientIPFromIncomingContext extracts a trusted client IP from gRPC metadata.
// 从 gRPC metadata 中提取可信的客户端 IP。
func GetClientIPFromIncomingContext(ctx context.Context) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	for _, key := range []string{"x-forwarded-for", "x-real-ip", "client-ip"} {
		values := md.Get(key)
		if len(values) == 0 {
			continue
		}

		candidate := strings.TrimSpace(values[0])
		if candidate == "" {
			continue
		}

		if key == "x-forwarded-for" {
			parts := strings.Split(candidate, ",")
			if len(parts) > 0 {
				return strings.TrimSpace(parts[0])
			}
		}

		return candidate
	}

	return ""
}
