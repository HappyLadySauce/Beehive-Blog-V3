package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
)

// TestRegisterServiceExecute verifies local registration behavior.
// TestRegisterServiceExecute 验证本地注册行为。
func TestRegisterServiceExecute(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)

	t.Run("register success", func(t *testing.T) {
		deps := newDeps(t, now)
		svc := service.NewRegisterService(deps)

		result, err := svc.Execute(context.Background(), service.RegisterLocalUserInput{
			Username: "alice_001",
			Email:    "alice@example.com",
			Password: "password123",
			Nickname: "Alice",
			ClientIP: "127.0.0.1",
		})
		if err != nil {
			t.Fatalf("expected register to succeed, got %v", err)
		}
		if result.User == nil || result.Session == nil {
			t.Fatalf("expected user and session to be returned")
		}
		if result.RefreshToken == "" || result.AccessToken == "" {
			t.Fatalf("expected token pair to be returned")
		}
	})

	t.Run("duplicate username", func(t *testing.T) {
		deps := newDeps(t, now)
		svc := service.NewRegisterService(deps)

		_, err := svc.Execute(context.Background(), service.RegisterLocalUserInput{
			Username: "alice_001",
			Email:    "alice@example.com",
			Password: "password123",
		})
		if err != nil {
			t.Fatalf("expected first register to succeed, got %v", err)
		}

		_, err = svc.Execute(context.Background(), service.RegisterLocalUserInput{
			Username: "alice_001",
			Email:    "alice2@example.com",
			Password: "password123",
		})
		if !service.IsKind(err, service.ErrorKindAlreadyExists) {
			t.Fatalf("expected already exists error, got %v", err)
		}
	})

	t.Run("invalid password", func(t *testing.T) {
		deps := newDeps(t, now)
		svc := service.NewRegisterService(deps)

		_, err := svc.Execute(context.Background(), service.RegisterLocalUserInput{
			Username: "alice_001",
			Email:    "alice@example.com",
			Password: "short",
		})
		if !service.IsKind(err, service.ErrorKindInvalidArgument) {
			t.Fatalf("expected invalid argument error, got %v", err)
		}
	})
}
