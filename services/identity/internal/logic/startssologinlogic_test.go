package logic

import (
	"context"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/config"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestStartSsoLoginRejectsProvidersThatAreNotReady verifies that phase-3 only exposes GitHub.
// TestStartSsoLoginRejectsProvidersThatAreNotReady 验证第三阶段只对外开放 GitHub。
func TestStartSsoLoginRejectsProvidersThatAreNotReady(t *testing.T) {
	t.Parallel()

	logic := NewStartSsoLoginLogic(context.Background(), &svc.ServiceContext{
		Config: config.Config{
			Security: config.SecurityConf{
				StateTTLSeconds: 600,
			},
			SSO: config.SSOConf{
				QQ: config.OAuthProviderConf{
					Enabled:     true,
					ClientID:    "qq-client-id",
					RedirectURL: "https://example.com/auth/qq/callback",
				},
			},
		},
	})

	_, err := logic.StartSsoLogin(&pb.StartSsoLoginRequest{
		Provider:    "qq",
		RedirectUri: "https://example.com/auth/qq/callback",
		State:       "fixed-state",
	})
	if err == nil {
		t.Fatalf("expected qq provider to be rejected")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected grpc status error, got %v", err)
	}
	if st.Code() != codes.FailedPrecondition {
		t.Fatalf("expected failed precondition, got %s", st.Code())
	}
	if st.Message() != "sso_provider_not_ready" {
		t.Fatalf("expected sso_provider_not_ready, got %s", st.Message())
	}
}
