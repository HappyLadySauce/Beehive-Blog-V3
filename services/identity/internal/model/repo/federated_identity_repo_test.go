package repo_test

import (
	"context"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/testkit"
)

// TestFederatedIdentityRepository verifies federated identity persistence.
// TestFederatedIdentityRepository 验证联邦身份持久化。
func TestFederatedIdentityRepository(t *testing.T) {
	store := testkit.NewStore(t)
	user := testkit.CreateUser(t, store)
	ctx := context.Background()

	fed := testkit.CreateFederatedIdentity(t, store, user.ID, "github", "12345")

	got, err := store.FederatedIdentities.GetByProviderSubject(ctx, "github", "12345")
	if err != nil {
		t.Fatalf("expected federated identity lookup to succeed, got %v", err)
	}
	if got.ID != fed.ID {
		t.Fatalf("expected federated identity id %d, got %d", fed.ID, got.ID)
	}

	t.Run("lookup by provider identity fields", func(t *testing.T) {
		openID := "wechat-openid-123"
		unionID := "wechat-unionid-456"
		wechatFed := testkit.CreateFederatedIdentity(t, store, user.ID, "wechat", unionID, func(fed *entity.FederatedIdentity) {
			fed.ProviderSubjectType = "unionid"
			fed.OpenID = &openID
			fed.UnionID = &unionID
		})

		found, err := store.FederatedIdentities.GetByProviderIdentity(ctx, "wechat", unionID, &openID, &unionID)
		if err != nil {
			t.Fatalf("expected provider identity lookup to succeed, got %v", err)
		}
		if found.ID != wechatFed.ID {
			t.Fatalf("expected federated identity id %d, got %d", wechatFed.ID, found.ID)
		}
	})

	t.Run("lookup by provider openid fallback", func(t *testing.T) {
		openID := "wechat-openid-789"
		unionID := "wechat-unionid-999"
		wechatFed := testkit.CreateFederatedIdentity(t, store, user.ID, "wechat", unionID, func(fed *entity.FederatedIdentity) {
			fed.ProviderSubjectType = "unionid"
			fed.OpenID = &openID
			fed.UnionID = &unionID
		})

		found, err := store.FederatedIdentities.GetByProviderIdentity(ctx, "wechat", "migrated-subject", &openID, nil)
		if err != nil {
			t.Fatalf("expected openid fallback lookup to succeed, got %v", err)
		}
		if found.ID != wechatFed.ID {
			t.Fatalf("expected federated identity id %d, got %d", wechatFed.ID, found.ID)
		}
	})

	lastLoginAt := time.Now().UTC()
	displayName := "Octocat"
	if err := store.FederatedIdentities.TouchLogin(ctx, fed.ID, nil, nil, nil, &displayName, nil, nil, nil, nil, []byte(`{"id":12345}`), lastLoginAt); err != nil {
		t.Fatalf("expected touch login to succeed, got %v", err)
	}
}
