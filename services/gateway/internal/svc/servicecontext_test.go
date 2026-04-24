package svc

import (
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/config"
	"github.com/zeromicro/go-zero/zrpc"
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

func TestNewServiceContextReturnsErrorOnInvalidContentRPCClient(t *testing.T) {
	t.Parallel()

	_, err := NewServiceContext(config.Config{
		IdentityRPC: config.IdentityRPCConf{
			RpcClientConf:      zrpc.NewDirectClientConf([]string{"127.0.0.1:1"}, "", ""),
			InternalAuthToken:  "secret",
			InternalCallerName: "gateway",
		},
		ContentRPC: config.ContentRPCConf{
			InternalAuthToken:  "secret",
			InternalCallerName: "gateway",
		},
		Security: config.GatewaySecurityConf{},
	})
	if err == nil {
		t.Fatalf("expected service context initialization to fail for invalid content rpc client config")
	}
}
