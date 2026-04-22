package provider_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth/provider"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/config"
)

// TestWeChatClientBuildAuthorizeURLUsesWebsiteScope verifies website login always uses snsapi_login.
// TestWeChatClientBuildAuthorizeURLUsesWebsiteScope 验证网站登录始终使用 snsapi_login。
func TestWeChatClientBuildAuthorizeURLUsesWebsiteScope(t *testing.T) {
	t.Parallel()

	client := provider.NewWeChatClient(config.OAuthProviderConf{
		Enabled:      true,
		ClientID:     "wechat-appid",
		ClientSecret: "wechat-secret",
		RedirectURL:  "https://example.com/auth/wechat/callback",
		Scopes:       []string{"snsapi_base"},
	})

	authURL, err := client.BuildAuthorizeURL("state-123")
	if err != nil {
		t.Fatalf("expected build authorize url to succeed, got %v", err)
	}

	trimmed := strings.TrimSuffix(authURL, "#wechat_redirect")
	parsed, err := url.Parse(trimmed)
	if err != nil {
		t.Fatalf("expected authorize url to be parseable, got %v", err)
	}
	if parsed.Query().Get("scope") != "snsapi_login" {
		t.Fatalf("expected website scope snsapi_login, got %q", parsed.Query().Get("scope"))
	}
}

// TestWeChatClientExchangeCode verifies WeChat OAuth code exchange.
// TestWeChatClientExchangeCode 验证微信 OAuth code 交换。
func TestWeChatClientExchangeCode(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/token" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"access_token":  "wechat-access-token",
			"refresh_token": "wechat-refresh-token",
			"openid":        "wechat-openid-123",
			"unionid":       "wechat-unionid-456",
			"scope":         "snsapi_login",
		})
	}))
	defer server.Close()

	client := provider.NewWeChatClient(config.OAuthProviderConf{
		Enabled:      true,
		ClientID:     "wechat-appid",
		ClientSecret: "wechat-secret",
		RedirectURL:  "https://example.com/auth/wechat/callback",
	})
	client.HTTPClient = server.Client()
	client.TokenURL = server.URL + "/token"

	token, err := client.ExchangeCode(context.Background(), "code-123", "https://example.com/auth/wechat/callback")
	if err != nil {
		t.Fatalf("expected exchange code to succeed, got %v", err)
	}
	if token == nil || token.Token != "wechat-access-token" || token.OpenID != "wechat-openid-123" || token.UnionID != "wechat-unionid-456" {
		t.Fatalf("expected complete wechat token context, got %#v", token)
	}
}

// TestWeChatClientFetchProfilePrefersUnionID verifies unionid-first normalization.
// TestWeChatClientFetchProfilePrefersUnionID 验证 unionid 优先归一化。
func TestWeChatClientFetchProfilePrefersUnionID(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/userinfo" {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"openid":     "wechat-openid-123",
			"unionid":    "wechat-unionid-456",
			"nickname":   "WeChat Nick",
			"headimgurl": "https://example.com/wechat-avatar.png",
		})
	}))
	defer server.Close()

	client := provider.NewWeChatClient(config.OAuthProviderConf{
		Enabled:      true,
		ClientID:     "wechat-appid",
		ClientSecret: "wechat-secret",
		RedirectURL:  "https://example.com/auth/wechat/callback",
		Scopes:       []string{"snsapi_login"},
	})
	client.HTTPClient = server.Client()
	client.UserInfoURL = server.URL + "/userinfo"

	profile, raw, err := client.FetchProfile(context.Background(), &provider.AccessToken{
		Token:   "wechat-access-token",
		OpenID:  "wechat-openid-123",
		UnionID: "wechat-unionid-456",
		Scope:   "snsapi_login",
	})
	if err != nil {
		t.Fatalf("expected fetch profile to succeed, got %v", err)
	}
	if profile.Subject != "wechat-unionid-456" || profile.SubjectType != "unionid" {
		t.Fatalf("expected unionid subject, got %#v", profile)
	}
	if profile.OpenID == nil || *profile.OpenID != "wechat-openid-123" {
		t.Fatalf("expected wechat openid to be retained, got %#v", profile.OpenID)
	}
	if len(raw) == 0 {
		t.Fatalf("expected raw profile payload to be present")
	}
}

// TestWeChatClientFetchProfileFallsBackToOpenID verifies openid fallback behavior.
// TestWeChatClientFetchProfileFallsBackToOpenID 验证 openid 回退行为。
func TestWeChatClientFetchProfileFallsBackToOpenID(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"openid":   "wechat-openid-123",
			"nickname": "WeChat Nick",
		})
	}))
	defer server.Close()

	client := provider.NewWeChatClient(config.OAuthProviderConf{Enabled: true, ClientID: "wechat-appid"})
	client.HTTPClient = server.Client()
	client.UserInfoURL = server.URL

	profile, _, err := client.FetchProfile(context.Background(), &provider.AccessToken{
		Token:  "wechat-access-token",
		OpenID: "wechat-openid-123",
	})
	if err != nil {
		t.Fatalf("expected fetch profile to succeed, got %v", err)
	}
	if profile.Subject != "wechat-openid-123" || profile.SubjectType != "openid" {
		t.Fatalf("expected openid fallback subject, got %#v", profile)
	}
}
