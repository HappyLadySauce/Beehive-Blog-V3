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
		// oauth2 treats "text/plain" token responses as form-encoded; force JSON so the library unmarshals wire JSON.
		// oauth2 将 text/plain 的 token 响应当作表单解析；强制 JSON，以便库按 JSON 解析。
		w.Header().Set("Content-Type", "application/json")
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
	// GitHub accepts client credentials in the POST body; pin AuthStyle to avoid
	// oauth2's two-step header/body probe issuing two token requests to the stub.
	// GitHub 接受将 client 凭证放在 POST body；固定 AuthStyle，避免 oauth2 对 stub 发起两次探测请求。
	client.OAuthEndpoint = oauth2.Endpoint{
		AuthURL:   server.URL + "/authorize",
		TokenURL:  server.URL + "/token",
		AuthStyle: oauth2.AuthStyleInParams,
	}

	token, err := client.ExchangeCode(context.Background(), "code-123", "https://example.com/callback")
	if err != nil {
		t.Fatalf("expected exchange code to succeed, got %v", err)
	}
	if token == nil || token.Token != "github-access-token" {
		t.Fatalf("expected github-access-token, got %#v", token)
	}
}

// TestGitHubClientFetchProfile verifies GitHub profile normalization.
// TestGitHubClientFetchProfile 验证 GitHub 用户资料规范化。
func TestGitHubClientFetchProfile(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/user":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":         12345,
				"login":      "octocat",
				"name":       "The Octocat",
				"email":      "OctoCat@Example.com",
				"avatar_url": "https://example.com/avatar.png",
			})
		case "/user/emails":
			_ = json.NewEncoder(w).Encode([]map[string]any{
				{
					"email":    "OctoCat@Example.com",
					"verified": true,
					"primary":  true,
				},
			})
		default:
			http.NotFound(w, r)
		}
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

	profile, raw, err := client.FetchProfile(context.Background(), &provider.AccessToken{Token: "access-token"})
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
	if !profile.EmailVerified {
		t.Fatalf("expected email to be marked verified")
	}
	if len(raw) == 0 {
		t.Fatalf("expected raw profile payload to be present")
	}
}

// TestGitHubClientFetchProfileWithoutVerifiedEmail verifies unverified email is not trusted.
// TestGitHubClientFetchProfileWithoutVerifiedEmail 验证未验证邮箱不会被信任。
func TestGitHubClientFetchProfileWithoutVerifiedEmail(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/user":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":    12345,
				"login": "octocat",
				"name":  "The Octocat",
				"email": "octocat@example.com",
			})
		case "/user/emails":
			_ = json.NewEncoder(w).Encode([]map[string]any{
				{
					"email":    "octocat@example.com",
					"verified": false,
					"primary":  true,
				},
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := provider.NewGitHubClient(config.OAuthProviderConf{
		Enabled:      true,
		ClientID:     " github-client-id ",
		ClientSecret: " github-client-secret ",
		RedirectURL:  " https://example.com/callback ",
		Scopes:       []string{" read:user ", " user:email "},
	})
	client.HTTPClient = server.Client()
	client.APIBaseURL = server.URL

	profile, _, err := client.FetchProfile(context.Background(), &provider.AccessToken{Token: "access-token"})
	if err != nil {
		t.Fatalf("expected fetch profile to succeed, got %v", err)
	}
	if profile.Email != nil {
		t.Fatalf("expected unverified email to be ignored, got %#v", profile.Email)
	}
	if profile.EmailVerified {
		t.Fatalf("expected email to be unverified")
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

	if _, _, err := client.FetchProfile(context.Background(), &provider.AccessToken{Token: "access-token"}); err == nil {
		t.Fatalf("expected incomplete profile to fail")
	}
}
