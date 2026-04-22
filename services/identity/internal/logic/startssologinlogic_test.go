package logic

import (
	"context"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/testkit"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
)

// TestStartSsoLoginAcceptsConfiguredProviders verifies transport forwarding for supported providers.
// TestStartSsoLoginAcceptsConfiguredProviders 验证已支持 provider 的 transport 转发。
func TestStartSsoLoginAcceptsConfiguredProviders(t *testing.T) {
	t.Parallel()

	svcCtx := testkit.NewServiceContext(t)
	logic := NewStartSsoLoginLogic(context.Background(), svcCtx)

	resp, err := logic.StartSsoLogin(&pb.StartSsoLoginRequest{
		Provider:    "qq",
		RedirectUri: svcCtx.Config.SSO.QQ.RedirectURL,
		State:       "fixed-state",
	})
	if err == nil {
		if resp.GetProvider() != "qq" || resp.GetAuthUrl() == "" {
			t.Fatalf("expected qq auth response, got %+v", resp)
		}
		return
	}

	t.Fatalf("expected qq provider to be supported, got %v", err)
}
