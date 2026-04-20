package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/testkit"
)

// TestCurrentUserServiceExecute verifies trusted current user lookup behavior.
// TestCurrentUserServiceExecute 验证可信当前用户查询行为。
func TestCurrentUserServiceExecute(t *testing.T) {
	now := time.Date(2026, 4, 20, 12, 0, 0, 0, time.UTC)

	t.Run("lookup success", func(t *testing.T) {
		deps := newDeps(t, now)
		user := testkit.CreateUser(t, deps.Store)

		svc := service.NewCurrentUserService(deps)
		result, err := svc.Execute(context.Background(), service.GetCurrentUserInput{UserID: user.ID})
		if err != nil {
			t.Fatalf("expected lookup to succeed, got %v", err)
		}
		if result.User.ID != user.ID {
			t.Fatalf("expected user_id=%d, got %d", user.ID, result.User.ID)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		deps := newDeps(t, now)
		svc := service.NewCurrentUserService(deps)

		_, err := svc.Execute(context.Background(), service.GetCurrentUserInput{UserID: 9999})
		if !service.IsKind(err, service.ErrorKindNotFound) {
			t.Fatalf("expected not found error, got %v", err)
		}
	})
}
