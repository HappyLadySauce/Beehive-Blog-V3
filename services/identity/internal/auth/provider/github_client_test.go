package provider_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth/provider"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/config"
	"golang.org/x/oauth2"
)

// TestGitHubClientExchangeCode verifies OAuth code exchange.
// TestGitHubClientExchangeCode 验证 OAuth code 交换流程。
func TestGitHubClientExchangeCode(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/token" {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token": "github-access-token",
			"token_type":   "bearer",
		})
	}))
	defer server.Close()

	client := provider.NewGitHubClient(config.OAuthProviderConf{
		Enabled:      true,
		ClientID:     "github-client-id",
		ClientSecret: "github-client-secret",
		RedirectURL:  "https://example.com/callback",
	})
	client.HTTPClient = server.Client()
	client.OAuthEndpoint = oauth2.Endpoint{
		AuthURL:  server.URL + "/authorize",
		TokenURL: server.URL + "/token",
	}

	token, err := client.ExchangeCode(context.Background(), "code-123", "https://example.com/callback")
	if err != nil {
		t.Fatalf("expected exchange code to succeed, got %v", err)
	}
	if token != "github-access-token" {
		t.Fatalf("expected github-access-token, got %q", token)
	}
}

// TestGitHubClientFetchProfile verifies GitHub profile normalization.
// TestGitHubClientFetchProfile 验证 GitHub 用户资料规范化。
func TestGitHubClientFetchProfile(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/user" {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id":         12345,
			"login":      "octocat",
			"name":       "The Octocat",
			"email":      "OctoCat@Example.com",
			"avatar_url": "https://example.com/avatar.png",
		})
	}))
	defer server.Close()

	client := provider.NewGitHubClient(config.OAuthProviderConf{
		Enabled:      true,
		ClientID:     "github-client-id",
		ClientSecret: "github-client-secret",
		RedirectURL:  "https://example.com/callback",
		Scopes:       []string{"read:user"},
	})
	client.HTTPClient = server.Client()
	client.APIBaseURL = server.URL

	profile, raw, err := client.FetchProfile(context.Background(), "access-token")
	if err != nil {
		t.Fatalf("expected fetch profile to succeed, got %v", err)
	}
	if profile.Subject != "12345" {
		t.Fatalf("expected subject 12345, got %q", profile.Subject)
	}
	if profile.Login != "octocat" {
		t.Fatalf("expected login octocat, got %q", profile.Login)
	}
	if profile.Email == nil || *profile.Email != "octocat@example.com" {
		t.Fatalf("expected normalized email, got %#v", profile.Email)
	}
	if len(raw) == 0 {
		t.Fatalf("expected raw profile payload to be present")
	}
}

// TestGitHubClientFetchProfileRejectsIncompleteProfile verifies required field checks.
// TestGitHubClientFetchProfileRejectsIncompleteProfile 验证缺少关键字段时的拒绝行为。
func TestGitHubClientFetchProfileRejectsIncompleteProfile(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": 0,
		})
	}))
	defer server.Close()

	client := provider.NewGitHubClient(config.OAuthProviderConf{Enabled: true})
	client.HTTPClient = server.Client()
	client.APIBaseURL = server.URL

	if _, _, err := client.FetchProfile(context.Background(), "access-token"); err == nil {
		t.Fatalf("expected incomplete profile to fail")
	}
}
