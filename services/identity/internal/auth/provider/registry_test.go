package provider_test

import (
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth/provider"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/config"
)

func readyConf(clientID, clientSecret, redirectURL string) config.OAuthProviderConf {
	return config.OAuthProviderConf{
		Enabled:      true,
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
	}
}

// TestRegistryGet verifies provider lookup behavior.
// TestRegistryGet 验证 provider 查询行为。
func TestRegistryGet(t *testing.T) {
	t.Parallel()

	registry := provider.NewRegistry(
		provider.NewGitHubClient(readyConf("github-client-id", "github-client-secret", "https://example.com/auth/github/callback")),
		provider.NewQQClient(readyConf("qq-client-id", "qq-client-secret", "https://example.com/auth/qq/callback")),
		provider.NewWeChatClient(readyConf("wechat-appid", "wechat-secret", "https://example.com/auth/wechat/callback")),
	)

	githubProvider, ok := registry.Get(auth.ProviderGitHub)
	if !ok {
		t.Fatalf("expected github provider to be registered")
	}
	if !githubProvider.Enabled() {
		t.Fatalf("expected github provider to be enabled")
	}
	if !githubProvider.LoginReady() {
		t.Fatalf("expected github provider to be login ready")
	}

	if _, ok := registry.Get(auth.ProviderQQ); !ok {
		t.Fatalf("expected qq provider to be registered")
	}
	if _, ok := registry.Get(auth.ProviderWeChat); !ok {
		t.Fatalf("expected wechat provider to be registered")
	}
}

// TestRegistryGetCallback verifies callback-capable provider lookup behavior.
// TestRegistryGetCallback 验证支持 callback 的 provider 查询行为。
func TestRegistryGetCallback(t *testing.T) {
	t.Parallel()

	registry := provider.NewRegistry(
		provider.NewGitHubClient(readyConf("github-client-id", "github-client-secret", "https://example.com/auth/github/callback")),
		provider.NewQQClient(readyConf("qq-client-id", "qq-client-secret", "https://example.com/auth/qq/callback")),
		provider.NewWeChatClient(readyConf("wechat-appid", "wechat-secret", "https://example.com/auth/wechat/callback")),
	)

	if _, ok := registry.GetCallback(auth.ProviderGitHub); !ok {
		t.Fatalf("expected github callback provider to be registered")
	}
	if _, ok := registry.GetCallback(auth.ProviderQQ); !ok {
		t.Fatalf("expected qq callback provider to be registered")
	}
	if _, ok := registry.GetCallback(auth.ProviderWeChat); !ok {
		t.Fatalf("expected wechat callback provider to be registered")
	}
}
