package service_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/testkit"
	"golang.org/x/oauth2"
)

// TestSSOFinishServiceExecute verifies GitHub callback completion behavior.
// TestSSOFinishServiceExecute 验证 GitHub callback 完成行为。
func TestSSOFinishServiceExecute(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)

	t.Run("github happy path", func(t *testing.T) {
		deps := newDeps(t, now)
		testkit.CreateOAuthState(t, deps.Store, auth.ProviderGitHub, "state-1", deps.Config.SSO.GitHub.RedirectURL)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/token":
				w.Header().Set("Content-Type", "application/json")
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
		if !errors.Is(err, errs.E(errs.CodeIdentitySSOProviderNotReady)) {
			t.Fatalf("expected provider not ready error, got %v", err)
		}
	})

	t.Run("invalid state fails", func(t *testing.T) {
		deps := newDeps(t, now)
		var tokenHits atomic.Int32

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/token":
				tokenHits.Add(1)
				w.Header().Set("Content-Type", "application/json")
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
		_, err := svc.Execute(context.Background(), service.FinishSSOInput{
			Provider:    auth.ProviderGitHub,
			Code:        "code-123",
			State:       "missing-state",
			RedirectURI: deps.Config.SSO.GitHub.RedirectURL,
		})
		if !errors.Is(err, errs.E(errs.CodeIdentitySSOStateInvalid)) {
			t.Fatalf("expected invalid sso state error, got %v", err)
		}
		if tokenHits.Load() != 0 {
			t.Fatalf("expected no upstream token exchange before invalid state rejection, got %d hits", tokenHits.Load())
		}
	})

	t.Run("consumed state fails", func(t *testing.T) {
		deps := newDeps(t, now)
		testkit.CreateOAuthState(t, deps.Store, auth.ProviderGitHub, "state-consumed", deps.Config.SSO.GitHub.RedirectURL, func(state *entity.OAuthLoginState) {
			consumedAt := now.Add(-time.Minute)
			state.ConsumedAt = &consumedAt
		})
		var tokenHits atomic.Int32

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/token":
				tokenHits.Add(1)
				w.Header().Set("Content-Type", "application/json")
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
		_, err := svc.Execute(context.Background(), service.FinishSSOInput{
			Provider:    auth.ProviderGitHub,
			Code:        "code-123",
			State:       "state-consumed",
			RedirectURI: deps.Config.SSO.GitHub.RedirectURL,
		})
		if !errors.Is(err, errs.E(errs.CodeIdentitySSOStateInvalid)) {
			t.Fatalf("expected invalid sso state error, got %v", err)
		}
		if tokenHits.Load() != 0 {
			t.Fatalf("expected no upstream token exchange for consumed state, got %d hits", tokenHits.Load())
		}
	})

	t.Run("expired state fails before upstream exchange", func(t *testing.T) {
		deps := newDeps(t, now)
		testkit.CreateOAuthState(t, deps.Store, auth.ProviderGitHub, "state-expired", deps.Config.SSO.GitHub.RedirectURL, func(state *entity.OAuthLoginState) {
			state.ExpiresAt = now.Add(-time.Minute)
		})
		var tokenHits atomic.Int32

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/token":
				tokenHits.Add(1)
				w.Header().Set("Content-Type", "application/json")
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
		_, err := svc.Execute(context.Background(), service.FinishSSOInput{
			Provider:    auth.ProviderGitHub,
			Code:        "code-123",
			State:       "state-expired",
			RedirectURI: deps.Config.SSO.GitHub.RedirectURL,
		})
		if !errors.Is(err, errs.E(errs.CodeIdentitySSOStateInvalid)) {
			t.Fatalf("expected invalid sso state error, got %v", err)
		}
		if tokenHits.Load() != 0 {
			t.Fatalf("expected no upstream token exchange for expired state, got %d hits", tokenHits.Load())
		}
	})

	t.Run("provider disabled fails", func(t *testing.T) {
		deps := newDeps(t, now)
		deps.Config.SSO.GitHub.Enabled = false
		deps.Providers = testkit.NewProviderRegistry(deps.Config)
		testkit.CreateOAuthState(t, deps.Store, auth.ProviderGitHub, "state-disabled", deps.Config.SSO.GitHub.RedirectURL)

		svc := service.NewSSOFinishService(deps)
		_, err := svc.Execute(context.Background(), service.FinishSSOInput{
			Provider:    auth.ProviderGitHub,
			Code:        "code-123",
			State:       "state-disabled",
			RedirectURI: deps.Config.SSO.GitHub.RedirectURL,
		})
		if !errors.Is(err, errs.E(errs.CodeIdentitySSOProviderDisabled)) {
			t.Fatalf("expected provider disabled error, got %v", err)
		}
	})

	t.Run("profile missing subject fails", func(t *testing.T) {
		deps := newDeps(t, now)
		testkit.CreateOAuthState(t, deps.Store, auth.ProviderGitHub, "state-bad-profile", deps.Config.SSO.GitHub.RedirectURL)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/token":
				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(map[string]any{
					"access_token": "github-access-token",
					"token_type":   "bearer",
				})
			case "/user":
				_ = json.NewEncoder(w).Encode(map[string]any{
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
		_, err := svc.Execute(context.Background(), service.FinishSSOInput{
			Provider:    auth.ProviderGitHub,
			Code:        "code-123",
			State:       "state-bad-profile",
			RedirectURI: deps.Config.SSO.GitHub.RedirectURL,
		})
		if !errors.Is(err, errs.E(errs.CodeIdentityInvalidCredentials)) {
			t.Fatalf("expected invalid credentials error, got %v", err)
		}
	})

	t.Run("exchange failure writes failure audit", func(t *testing.T) {
		deps := newDeps(t, now)
		testkit.CreateOAuthState(t, deps.Store, auth.ProviderGitHub, "state-exchange-failed", deps.Config.SSO.GitHub.RedirectURL)

		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/token":
				http.Error(w, "upstream failure", http.StatusBadGateway)
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
		_, err := svc.Execute(context.Background(), service.FinishSSOInput{
			Provider:    auth.ProviderGitHub,
			Code:        "code-exchange-failed",
			State:       "state-exchange-failed",
			RedirectURI: deps.Config.SSO.GitHub.RedirectURL,
			ClientIP:    "127.0.0.1",
			UserAgent:   "go-test",
		})
		if !errors.Is(err, errs.E(errs.CodeIdentityInvalidCredentials)) {
			t.Fatalf("expected invalid credentials error, got %v", err)
		}

		var audits []entity.IdentityAudit
		if err := deps.Store.DB().
			WithContext(context.Background()).
			Where("event_type = ? AND result = ?", auth.AuditEventFinishSSO, auth.AuditResultFailure).
			Find(&audits).Error; err != nil {
			t.Fatalf("failed to load audits: %v", err)
		}
		if len(audits) == 0 {
			t.Fatalf("expected at least one failure audit")
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
				w.Header().Set("Content-Type", "application/json")
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
