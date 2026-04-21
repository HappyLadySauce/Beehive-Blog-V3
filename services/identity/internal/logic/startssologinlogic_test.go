package logic

import (
	"context"
	"testing"

	identityprovider "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth/provider"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/config"
	identityservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestStartSsoLoginRejectsProvidersThatAreNotReady verifies that phase-3 only exposes GitHub.
// TestStartSsoLoginRejectsProvidersThatAreNotReady 验证第三阶段只对外开放 GitHub。
func TestStartSsoLoginRejectsProvidersThatAreNotReady(t *testing.T) {
	t.Parallel()

	conf := config.Config{
		Security: config.SecurityConf{
			StateTTLSeconds: 600,
		},
		SSO: config.SSOConf{
			GitHub: config.OAuthProviderConf{
				Enabled:      true,
				ClientID:     "github-client-id",
				ClientSecret: "github-client-secret",
				RedirectURL:  "https://example.com/auth/github/callback",
			},
		},
	}
	providers := identityprovider.NewRegistry(
		identityprovider.NewGitHubClient(conf.SSO.GitHub),
	)

	logic := NewStartSsoLoginLogic(context.Background(), &svc.ServiceContext{
		Config:    conf,
		Providers: providers,
		Services:  identityservice.NewManager(conf, nil, providers, nil),
	})

	_, err := logic.StartSsoLogin(&pb.StartSsoLoginRequest{
		Provider:    "qq",
		RedirectUri: "https://example.com/auth/github/callback",
		State:       "fixed-state",
	})
	if err == nil {
		t.Fatalf("expected qq provider to be rejected")
	}

	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected grpc status error, got %v", err)
	}
	if st.Code() != codes.InvalidArgument {
		t.Fatalf("expected invalid argument, got %s", st.Code())
	}
	if st.Message() != "unsupported provider" {
		t.Fatalf("expected unsupported provider, got %s", st.Message())
	}
}
