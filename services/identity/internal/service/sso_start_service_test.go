package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
)

// TestSSOStartServiceExecute verifies outbound SSO authorize behavior.
// TestSSOStartServiceExecute 验证对外 SSO 授权行为。
func TestSSOStartServiceExecute(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)

	t.Run("github success", func(t *testing.T) {
		deps := newDeps(t, now)
		svc := service.NewSSOStartService(deps)

		result, err := svc.Execute(context.Background(), service.StartSSOInput{
			Provider:    auth.ProviderGitHub,
			RedirectURI: deps.Config.SSO.GitHub.RedirectURL,
			State:       "fixed-state",
			ClientIP:    "127.0.0.1",
		})
		if err != nil {
			t.Fatalf("expected github sso start to succeed, got %v", err)
		}
		if result.AuthURL == "" {
			t.Fatalf("expected auth url to be non-empty")
		}
	})

	t.Run("qq rejected", func(t *testing.T) {
		deps := newDeps(t, now)
		svc := service.NewSSOStartService(deps)

		_, err := svc.Execute(context.Background(), service.StartSSOInput{
			Provider:    auth.ProviderQQ,
			RedirectURI: deps.Config.SSO.QQ.RedirectURL,
			State:       "fixed-state",
		})
		if !errors.Is(err, errs.E(errs.CodeIdentitySSOProviderNotReady)) {
			t.Fatalf("expected provider not ready error, got %v", err)
		}
	})

	t.Run("redirect uri mismatch", func(t *testing.T) {
		deps := newDeps(t, now)
		svc := service.NewSSOStartService(deps)

		_, err := svc.Execute(context.Background(), service.StartSSOInput{
			Provider:    auth.ProviderGitHub,
			RedirectURI: "https://example.com/other/callback",
			State:       "fixed-state",
		})
		if !errors.Is(err, errs.E(errs.CodeIdentityInvalidArgument)) {
			t.Fatalf("expected invalid argument error, got %v", err)
		}
	})
}
