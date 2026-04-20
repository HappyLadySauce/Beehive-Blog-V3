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
	"golang.org/x/oauth2"
)

// TestSSOFinishServiceExecute verifies GitHub callback completion behavior.
// TestSSOFinishServiceExecute 验证 GitHub callback 完成行为。
func TestSSOFinishServiceExecute(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)

	t.Run("github happy path", func(t *testing.T) {
		deps := newDeps(t, now)
		testkit.CreateOAuthState(t, deps.Store, auth.ProviderGitHub, "state-1", deps.Config.SSO.GitHub.RedirectURL)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/token":
				_ = json.NewEncoder(w).Encode(map[string]any{
					"access_token": "github-access-token",
					"token_type":   "bearer",
				})
			case "/user":
				_ = json.NewEncoder(w).Encode(map[string]any{
					"id":         12345,
					"login":      "octocat",
					"name":       "The Octocat",
					"email":      "octocat@example.com",
					"avatar_url": "https://example.com/avatar.png",
				})
			default:
				http.NotFound(w, r)
			}
		}))
		defer server.Close()

		client := githubClientFromDeps(t, deps)
		client.HTTPClient = server.Client()
		client.APIBaseURL = server.URL
		client.OAuthEndpoint = oauth2.Endpoint{
			AuthURL:  server.URL + "/authorize",
			TokenURL: server.URL + "/token",
		}

		svc := service.NewSSOFinishService(deps)
		result, err := svc.Execute(context.Background(), service.FinishSSOInput{
			Provider:    auth.ProviderGitHub,
			Code:        "code-123",
			State:       "state-1",
			RedirectURI: deps.Config.SSO.GitHub.RedirectURL,
			ClientIP:    "127.0.0.1",
		})
		if err != nil {
			t.Fatalf("expected github sso finish to succeed, got %v", err)
		}
		if result.User == nil || result.Session == nil {
			t.Fatalf("expected user and session to be returned")
		}
	})

	t.Run("provider not ready", func(t *testing.T) {
		deps := newDeps(t, now)
		svc := service.NewSSOFinishService(deps)

		_, err := svc.Execute(context.Background(), service.FinishSSOInput{
			Provider:    auth.ProviderQQ,
			Code:        "code-123",
			State:       "state-1",
			RedirectURI: deps.Config.SSO.QQ.RedirectURL,
		})
		if !service.IsKind(err, service.ErrorKindUnimplemented) {
			t.Fatalf("expected unimplemented error, got %v", err)
		}
	})

	t.Run("existing federated identity reuses user", func(t *testing.T) {
		deps := newDeps(t, now)
		user := testkit.CreateUser(t, deps.Store)
		testkit.CreateFederatedIdentity(t, deps.Store, user.ID, auth.ProviderGitHub, "12345")
		testkit.CreateOAuthState(t, deps.Store, auth.ProviderGitHub, "state-2", deps.Config.SSO.GitHub.RedirectURL)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/token":
				_ = json.NewEncoder(w).Encode(map[string]any{
					"access_token": "github-access-token",
					"token_type":   "bearer",
				})
			case "/user":
				_ = json.NewEncoder(w).Encode(map[string]any{
					"id":    12345,
					"login": "octocat",
					"name":  "The Octocat",
				})
			default:
				http.NotFound(w, r)
			}
		}))
		defer server.Close()

		client := githubClientFromDeps(t, deps)
		client.HTTPClient = server.Client()
		client.APIBaseURL = server.URL
		client.OAuthEndpoint = oauth2.Endpoint{
			AuthURL:  server.URL + "/authorize",
			TokenURL: server.URL + "/token",
		}

		svc := service.NewSSOFinishService(deps)
		result, err := svc.Execute(context.Background(), service.FinishSSOInput{
			Provider:    auth.ProviderGitHub,
			Code:        "code-456",
			State:       "state-2",
			RedirectURI: deps.Config.SSO.GitHub.RedirectURL,
		})
		if err != nil {
			t.Fatalf("expected github sso finish to succeed, got %v", err)
		}
		if result.User.ID != user.ID {
			t.Fatalf("expected existing user_id=%d, got %d", user.ID, result.User.ID)
		}
	})
}
