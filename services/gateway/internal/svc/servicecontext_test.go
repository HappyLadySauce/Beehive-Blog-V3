package svc

import (
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
)

// TestNewServiceContextReturnsErrorOnInvalidRPCClient verifies startup failures return errors instead of panicking.
// TestNewServiceContextReturnsErrorOnInvalidRPCClient 验证 RPC client 初始化失败时返回 error 而不是 panic。
func TestNewServiceContextReturnsErrorOnInvalidRPCClient(t *testing.T) {
	t.Parallel()

	_, err := NewServiceContext(config.Config{
		IdentityRPC: config.IdentityRPCConf{
			InternalAuthToken:  "secret",
			InternalCallerName: "gateway",
		},
		Security: config.GatewaySecurityConf{},
	})
	if err == nil {
		t.Fatalf("expected service context initialization to fail for invalid rpc client config")
	}
}
