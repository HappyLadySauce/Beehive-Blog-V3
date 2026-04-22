package provider_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth/provider"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/config"
)

// TestQQClientExchangeCode verifies QQ OAuth code exchange and openid retrieval.
// TestQQClientExchangeCode 验证 QQ OAuth code 交换与 openid 获取。
func TestQQClientExchangeCode(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/token":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token": "qq-access-token",
				"expires_in":   7776000,
			})
		case "/me":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"client_id": "qq-client-id",
				"openid":    "qq-openid-123",
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := provider.NewQQClient(config.OAuthProviderConf{
		Enabled:      true,
		ClientID:     "qq-client-id",
		ClientSecret: "qq-client-secret",
		RedirectURL:  "https://example.com/auth/qq/callback",
		Scopes:       []string{"get_user_info"},
	})
	client.HTTPClient = server.Client()
	client.TokenURL = server.URL + "/token"
	client.OpenIDURL = server.URL + "/me"

	token, err := client.ExchangeCode(context.Background(), "code-123", "https://example.com/auth/qq/callback")
	if err != nil {
		t.Fatalf("expected exchange code to succeed, got %v", err)
	}
	if token == nil || token.Token != "qq-access-token" || token.OpenID != "qq-openid-123" {
		t.Fatalf("expected qq access token with openid, got %#v", token)
	}
}

// TestQQClientExchangeCodeAcceptsStringErrorCode verifies tolerant QQ token error parsing.
// TestQQClientExchangeCodeAcceptsStringErrorCode 验证 QQ token 字符串错误码的宽容解析。
func TestQQClientExchangeCodeAcceptsStringErrorCode(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/token":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"error":"100013","error_description":"invalid code"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := provider.NewQQClient(config.OAuthProviderConf{
		Enabled:      true,
		ClientID:     "qq-client-id",
		ClientSecret: "qq-client-secret",
		RedirectURL:  "https://example.com/auth/qq/callback",
	})
	client.HTTPClient = server.Client()
	client.TokenURL = server.URL + "/token"
	client.OpenIDURL = server.URL + "/me"

	if _, err := client.ExchangeCode(context.Background(), "code-123", "https://example.com/auth/qq/callback"); err == nil {
		t.Fatalf("expected qq token string error code to fail")
	}
}

// TestQQClientFetchProfile verifies QQ profile normalization.
// TestQQClientFetchProfile 验证 QQ 用户资料规范化。
func TestQQClientFetchProfile(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/userinfo" {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{
			"ret":            0,
			"nickname":       "QQ Nick",
			"figureurl_qq_2": "https://example.com/qq-avatar.png",
		})
	}))
	defer server.Close()

	client := provider.NewQQClient(config.OAuthProviderConf{
		Enabled:      true,
		ClientID:     "qq-client-id",
		ClientSecret: "qq-client-secret",
		RedirectURL:  "https://example.com/auth/qq/callback",
		Scopes:       []string{"get_user_info"},
	})
	client.HTTPClient = server.Client()
	client.UserInfoURL = server.URL + "/userinfo"

	profile, raw, err := client.FetchProfile(context.Background(), &provider.AccessToken{
		Token:  "qq-access-token",
		OpenID: "qq-openid-123",
		Scope:  "get_user_info",
	})
	if err != nil {
		t.Fatalf("expected fetch profile to succeed, got %v", err)
	}
	if profile.Subject != "qq-openid-123" || profile.SubjectType != "openid" {
		t.Fatalf("expected qq openid subject, got %#v", profile)
	}
	if profile.DisplayName != "QQ Nick" {
		t.Fatalf("expected nickname to be normalized, got %#v", profile)
	}
	if profile.AvatarURL == nil || *profile.AvatarURL != "https://example.com/qq-avatar.png" {
		t.Fatalf("expected qq avatar url, got %#v", profile.AvatarURL)
	}
	if len(raw) == 0 {
		t.Fatalf("expected raw profile payload to be present")
	}
}

// TestQQClientFetchProfileRejectsIncompleteProfile verifies required QQ profile fields.
// TestQQClientFetchProfileRejectsIncompleteProfile 验证 QQ 资料缺少关键字段时的拒绝行为。
func TestQQClientFetchProfileRejectsIncompleteProfile(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"ret": 0,
		})
	}))
	defer server.Close()

	client := provider.NewQQClient(config.OAuthProviderConf{Enabled: true, ClientID: "qq-client-id"})
	client.HTTPClient = server.Client()
	client.UserInfoURL = server.URL

	if _, _, err := client.FetchProfile(context.Background(), &provider.AccessToken{
		Token:  "qq-access-token",
		OpenID: "qq-openid-123",
	}); err == nil {
		t.Fatalf("expected incomplete profile to fail")
	}
}

// TestQQClientFetchProfileAcceptsStringRet verifies tolerant QQ ret parsing.
// TestQQClientFetchProfileAcceptsStringRet 验证 QQ ret 字段字符串解析的宽容性。
func TestQQClientFetchProfileAcceptsStringRet(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"ret":"0","nickname":"QQ Nick","figureurl_qq_2":"https://example.com/qq-avatar.png"}`))
	}))
	defer server.Close()

	client := provider.NewQQClient(config.OAuthProviderConf{
		Enabled:      true,
		ClientID:     "qq-client-id",
		ClientSecret: "qq-client-secret",
		RedirectURL:  "https://example.com/auth/qq/callback",
	})
	client.HTTPClient = server.Client()
	client.UserInfoURL = server.URL

	profile, _, err := client.FetchProfile(context.Background(), &provider.AccessToken{
		Token:  "qq-access-token",
		OpenID: "qq-openid-123",
	})
	if err != nil {
		t.Fatalf("expected string ret profile fetch to succeed, got %v", err)
	}
	if profile == nil || profile.Subject != "qq-openid-123" {
		t.Fatalf("expected qq profile to be returned, got %#v", profile)
	}
}

// TestQQClientFetchProfileUsesStructuredErrorMessage verifies userinfo errors prefer structured messages.
// TestQQClientFetchProfileUsesStructuredErrorMessage 验证 userinfo 错误优先使用结构化错误信息。
func TestQQClientFetchProfileUsesStructuredErrorMessage(t *testing.T) {
	t.Parallel()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"ret":10001,"msg":"nickname denied"}`))
	}))
	defer server.Close()

	client := provider.NewQQClient(config.OAuthProviderConf{
		Enabled:      true,
		ClientID:     "qq-client-id",
		ClientSecret: "qq-client-secret",
		RedirectURL:  "https://example.com/auth/qq/callback",
	})
	client.HTTPClient = server.Client()
	client.UserInfoURL = server.URL

	_, _, err := client.FetchProfile(context.Background(), &provider.AccessToken{
		Token:  "qq-access-token",
		OpenID: "qq-openid-123",
	})
	if err == nil {
		t.Fatalf("expected structured qq userinfo error")
	}
	if got := err.Error(); got != "qq user info API returned error 10001: nickname denied" {
		t.Fatalf("expected structured qq error message, got %q", got)
	}
}
