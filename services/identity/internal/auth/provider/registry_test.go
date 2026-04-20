package provider_test

import (
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth/provider"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/config"
)

// TestRegistryGet verifies provider lookup behavior.
// TestRegistryGet 验证 provider 查询行为。
func TestRegistryGet(t *testing.T) {
	t.Parallel()

	registry := provider.NewRegistry(
		provider.NewGitHubClient(config.OAuthProviderConf{Enabled: true}),
		provider.NewQQClient(config.OAuthProviderConf{Enabled: true}),
		provider.NewWeChatClient(config.OAuthProviderConf{Enabled: true}),
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

	qqProvider, ok := registry.Get(auth.ProviderQQ)
	if !ok {
		t.Fatalf("expected qq provider to be registered")
	}
	if !qqProvider.Enabled() {
		t.Fatalf("expected qq provider to be enabled")
	}
	if qqProvider.LoginReady() {
		t.Fatalf("expected qq provider to be not login ready")
	}
}

// TestRegistryGetCallback verifies callback-capable provider lookup behavior.
// TestRegistryGetCallback 验证支持 callback 的 provider 查询行为。
func TestRegistryGetCallback(t *testing.T) {
	t.Parallel()

	registry := provider.NewRegistry(
		provider.NewGitHubClient(config.OAuthProviderConf{Enabled: true}),
		provider.NewQQClient(config.OAuthProviderConf{Enabled: true}),
	)

	if _, ok := registry.GetCallback(auth.ProviderGitHub); !ok {
		t.Fatalf("expected github callback provider to be registered")
	}
	if _, ok := registry.GetCallback(auth.ProviderQQ); ok {
		t.Fatalf("expected qq callback provider to be unavailable")
	}
}
