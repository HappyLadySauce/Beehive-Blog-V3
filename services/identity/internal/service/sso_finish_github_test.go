package service_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/testkit"
	"golang.org/x/oauth2"
)

func TestSSOFinishGitHubHappyPath(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)

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
			writeGitHubUserResponse(w, 12345, "octocat", "The Octocat", "octocat@example.com")
		case "/user/emails":
			writeGitHubEmailsResponse(w, []map[string]any{
				{
					"email":    "octocat@example.com",
					"verified": true,
					"primary":  true,
				},
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
}

func TestSSOFinishExistingFederatedIdentityReusesUser(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)

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
			writeGitHubUserResponse(w, 12345, "octocat", "The Octocat", "")
		case "/user/emails":
			writeGitHubEmailsResponse(w, nil)
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
}

func TestSSOFinishExistingFederatedIdentityIgnoresGitHubEmailsEndpointFailure(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)

	deps := newDeps(t, now)
	user := testkit.CreateUser(t, deps.Store)
	testkit.CreateFederatedIdentity(t, deps.Store, user.ID, auth.ProviderGitHub, "12345")
	testkit.CreateOAuthState(t, deps.Store, auth.ProviderGitHub, "state-emails-fallback", deps.Config.SSO.GitHub.RedirectURL)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/token":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token": "github-access-token",
				"token_type":   "bearer",
			})
		case "/user":
			writeGitHubUserResponse(w, 12345, "octocat", "The Octocat", "")
		case "/user/emails":
			http.Error(w, "forbidden", http.StatusForbidden)
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
		Code:        "code-emails-fallback",
		State:       "state-emails-fallback",
		RedirectURI: deps.Config.SSO.GitHub.RedirectURL,
	})
	if err != nil {
		t.Fatalf("expected federated identity login to survive emails fallback, got %v", err)
	}
	if result.User == nil || result.User.ID != user.ID {
		t.Fatalf("expected existing user_id=%d, got %#v", user.ID, result.User)
	}
}

func TestSSOFinishVerifiedGitHubEmailLinksExistingLocalUser(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)

	deps := newDeps(t, now)
	existingUser := testkit.CreateUser(t, deps.Store, func(user *entity.User) {
		email := "octocat@example.com"
		user.Email = &email
	})
	testkit.CreateOAuthState(t, deps.Store, auth.ProviderGitHub, "state-link-verified", deps.Config.SSO.GitHub.RedirectURL)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/token":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token": "github-access-token",
				"token_type":   "bearer",
			})
		case "/user":
			writeGitHubUserResponse(w, 99999, "octocat", "The Octocat", "octocat@example.com")
		case "/user/emails":
			writeGitHubEmailsResponse(w, []map[string]any{
				{
					"email":    "octocat@example.com",
					"verified": true,
					"primary":  true,
				},
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
		Code:        "code-link-verified",
		State:       "state-link-verified",
		RedirectURI: deps.Config.SSO.GitHub.RedirectURL,
	})
	if err != nil {
		t.Fatalf("expected verified github email linking to succeed, got %v", err)
	}
	if result.User == nil || result.User.ID != existingUser.ID {
		t.Fatalf("expected existing user_id=%d, got %#v", existingUser.ID, result.User)
	}
}

func TestSSOFinishUnverifiedGitHubEmailDoesNotLinkExistingLocalUser(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)

	deps := newDeps(t, now)
	existingUser := testkit.CreateUser(t, deps.Store, func(user *entity.User) {
		email := "octocat@example.com"
		user.Email = &email
	})
	testkit.CreateOAuthState(t, deps.Store, auth.ProviderGitHub, "state-link-unverified", deps.Config.SSO.GitHub.RedirectURL)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/token":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token": "github-access-token",
				"token_type":   "bearer",
			})
		case "/user":
			writeGitHubUserResponse(w, 88888, "octocat-unverified", "The Octocat", "octocat@example.com")
		case "/user/emails":
			writeGitHubEmailsResponse(w, []map[string]any{
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
		Code:        "code-link-unverified",
		State:       "state-link-unverified",
		RedirectURI: deps.Config.SSO.GitHub.RedirectURL,
	})
	if err != nil {
		t.Fatalf("expected unverified github email flow to succeed, got %v", err)
	}
	if result.User == nil {
		t.Fatalf("expected user to be returned")
	}
	if result.User.ID == existingUser.ID {
		t.Fatalf("expected unverified github email not to link existing user_id=%d", existingUser.ID)
	}
}

func TestSSOFinishGitHubEmailsEndpointFailureDoesNotAutoLinkExistingLocalEmailUser(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)

	deps := newDeps(t, now)
	existingUser := testkit.CreateUser(t, deps.Store, func(user *entity.User) {
		email := "octocat@example.com"
		user.Email = &email
	})
	testkit.CreateOAuthState(t, deps.Store, auth.ProviderGitHub, "state-link-fallback", deps.Config.SSO.GitHub.RedirectURL)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/token":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"access_token": "github-access-token",
				"token_type":   "bearer",
			})
		case "/user":
			writeGitHubUserResponse(w, 77777, "octocat-fallback", "The Octocat", "octocat@example.com")
		case "/user/emails":
			http.Error(w, "forbidden", http.StatusForbidden)
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
		Code:        "code-link-fallback",
		State:       "state-link-fallback",
		RedirectURI: deps.Config.SSO.GitHub.RedirectURL,
	})
	if err != nil {
		t.Fatalf("expected emails fallback flow to succeed, got %v", err)
	}
	if result.User == nil {
		t.Fatalf("expected user to be returned")
	}
	if result.User.ID == existingUser.ID {
		t.Fatalf("expected emails fallback not to auto link existing user_id=%d", existingUser.ID)
	}
}
