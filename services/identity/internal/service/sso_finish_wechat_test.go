package service_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/testkit"
)

func TestSSOFinishWeChatHappyPathPrefersUnionID(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)

	deps := newDeps(t, now)
	testkit.CreateOAuthState(t, deps.Store, auth.ProviderWeChat, "wechat-state-1", deps.Config.SSO.WeChat.RedirectURL)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/token":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token":  "wechat-access-token",
				"refresh_token": "wechat-refresh-token",
				"openid":        "wechat-openid-123",
				"unionid":       "wechat-unionid-456",
				"scope":         "snsapi_login",
			})
		case "/userinfo":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"openid":     "wechat-openid-123",
				"unionid":    "wechat-unionid-456",
				"nickname":   "WeChat Nick",
				"headimgurl": "https://example.com/wechat-avatar.png",
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := wechatClientFromDeps(t, deps)
	client.HTTPClient = server.Client()
	client.AuthorizeURL = server.URL + "/authorize"
	client.TokenURL = server.URL + "/token"
	client.UserInfoURL = server.URL + "/userinfo"

	svc := service.NewSSOFinishService(deps)
	result, err := svc.Execute(context.Background(), service.FinishSSOInput{
		Provider:    auth.ProviderWeChat,
		Code:        "code-wechat-123",
		State:       "wechat-state-1",
		RedirectURI: deps.Config.SSO.WeChat.RedirectURL,
		ClientIP:    "127.0.0.1",
	})
	if err != nil {
		t.Fatalf("expected wechat sso finish to succeed, got %v", err)
	}
	if result.User == nil || result.Session == nil {
		t.Fatalf("expected user and session to be returned")
	}
}
