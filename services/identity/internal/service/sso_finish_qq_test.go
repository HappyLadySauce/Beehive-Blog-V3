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

func TestSSOFinishQQHappyPath(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)

	deps := newDeps(t, now)
	testkit.CreateOAuthState(t, deps.Store, auth.ProviderQQ, "qq-state-1", deps.Config.SSO.QQ.RedirectURL)

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
		case "/userinfo":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"ret":            0,
				"nickname":       "QQ Nick",
				"figureurl_qq_2": "https://example.com/qq-avatar.png",
			})
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	client := qqClientFromDeps(t, deps)
	client.HTTPClient = server.Client()
	client.AuthorizeURL = server.URL + "/authorize"
	client.TokenURL = server.URL + "/token"
	client.OpenIDURL = server.URL + "/me"
	client.UserInfoURL = server.URL + "/userinfo"

	svc := service.NewSSOFinishService(deps)
	result, err := svc.Execute(context.Background(), service.FinishSSOInput{
		Provider:    auth.ProviderQQ,
		Code:        "code-qq-123",
		State:       "qq-state-1",
		RedirectURI: deps.Config.SSO.QQ.RedirectURL,
		ClientIP:    "127.0.0.1",
	})
	if err != nil {
		t.Fatalf("expected qq sso finish to succeed, got %v", err)
	}
	if result.User == nil || result.Session == nil {
		t.Fatalf("expected user and session to be returned")
	}
}
