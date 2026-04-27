package service_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/testkit"
)

// TestUserManagementServiceListUsers verifies admin user listing and authorization.
// TestUserManagementServiceListUsers 验证管理员用户列表与授权。
func TestUserManagementServiceListUsers(t *testing.T) {
	now := time.Date(2026, 4, 27, 10, 0, 0, 0, time.UTC)
	deps := newDeps(t, now)
	admin := createUserWithRole(t, deps.Store, auth.UserRoleAdmin, auth.UserStatusActive)
	member := createUserWithRole(t, deps.Store, auth.UserRoleMember, auth.UserStatusActive)
	createUserWithRole(t, deps.Store, auth.UserRoleMember, auth.UserStatusDisabled)

	svc := service.NewUserManagementService(deps)
	result, err := svc.ListUsers(context.Background(), service.ListUsersInput{
		ActorUserID: admin.ID,
		Role:        auth.UserRoleMember,
		Status:      auth.UserStatusActive,
		Page:        1,
		PageSize:    10,
	})
	if err != nil {
		t.Fatalf("expected list users to succeed, got %v", err)
	}
	if result.Total != 1 || len(result.Items) != 1 || result.Items[0].ID != member.ID {
		t.Fatalf("expected only active member, got total=%d items=%v", result.Total, result.Items)
	}

	_, err = svc.ListUsers(context.Background(), service.ListUsersInput{ActorUserID: member.ID})
	if !errors.Is(err, errs.E(errs.CodeIdentityAccessForbidden)) {
		t.Fatalf("expected access forbidden for non-admin, got %v", err)
	}
}

// TestUserManagementServiceProfileAndPassword verifies self-service mutations.
// TestUserManagementServiceProfileAndPassword 验证用户自助资料和密码修改。
func TestUserManagementServiceProfileAndPassword(t *testing.T) {
	now := time.Date(2026, 4, 27, 11, 0, 0, 0, time.UTC)
	deps := newDeps(t, now)
	user := createUserWithRole(t, deps.Store, auth.UserRoleMember, auth.UserStatusActive)
	oldHash, err := auth.HashPassword("OldPass123!", deps.Config.Security.PasswordHashCost)
	if err != nil {
		t.Fatalf("failed to hash old password: %v", err)
	}
	testkit.CreateCredentialLocal(t, deps.Store, user.ID, oldHash)

	svc := service.NewUserManagementService(deps)
	profile, err := svc.UpdateOwnProfile(context.Background(), service.UpdateOwnProfileInput{
		UserID:    user.ID,
		Nickname:  stringPtr("Alice"),
		AvatarURL: stringPtr("https://cdn.example.com/avatar.png"),
	})
	if err != nil {
		t.Fatalf("expected profile update to succeed, got %v", err)
	}
	if profile.User.Nickname == nil || *profile.User.Nickname != "Alice" {
		t.Fatalf("expected nickname to update, got %#v", profile.User.Nickname)
	}
	profile, err = svc.UpdateOwnProfile(context.Background(), service.UpdateOwnProfileInput{
		UserID:    user.ID,
		AvatarURL: stringPtr(""),
	})
	if err != nil {
		t.Fatalf("expected avatar-only profile update to succeed, got %v", err)
	}
	if profile.User.Nickname == nil || *profile.User.Nickname != "Alice" {
		t.Fatalf("expected omitted nickname to be preserved, got %#v", profile.User.Nickname)
	}
	if profile.User.AvatarURL != nil {
		t.Fatalf("expected explicit empty avatar URL to clear avatar, got %#v", profile.User.AvatarURL)
	}

	if err := svc.ChangeOwnPassword(context.Background(), service.ChangeOwnPasswordInput{
		UserID:      user.ID,
		OldPassword: "OldPass123!",
		NewPassword: "NewPass123!",
	}); err != nil {
		t.Fatalf("expected password change to succeed, got %v", err)
	}
	credential, err := deps.Store.CredentialLocals.GetByUserID(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("expected credential lookup to succeed, got %v", err)
	}
	if err := auth.VerifyPassword(credential.PasswordHash, "NewPass123!"); err != nil {
		t.Fatalf("expected new password to verify, got %v", err)
	}
}

// TestUserManagementServiceAdminMutationsAndAudits verifies admin writes and audit queries.
// TestUserManagementServiceAdminMutationsAndAudits 验证管理员写入与审计查询。
func TestUserManagementServiceAdminMutationsAndAudits(t *testing.T) {
	now := time.Date(2026, 4, 27, 12, 0, 0, 0, time.UTC)
	deps := newDeps(t, now)
	admin := createUserWithRole(t, deps.Store, auth.UserRoleAdmin, auth.UserStatusActive)
	target := createUserWithRole(t, deps.Store, auth.UserRoleMember, auth.UserStatusActive)

	svc := service.NewUserManagementService(deps)
	roleResult, err := svc.UpdateUserRole(context.Background(), service.UpdateUserRoleInput{
		ActorUserID:  admin.ID,
		TargetUserID: target.ID,
		Role:         auth.UserRoleAdmin,
	})
	if err != nil {
		t.Fatalf("expected role update to succeed, got %v", err)
	}
	if roleResult.User.Role != auth.UserRoleAdmin {
		t.Fatalf("expected role admin, got %s", roleResult.User.Role)
	}

	statusResult, err := svc.UpdateUserStatus(context.Background(), service.UpdateUserStatusInput{
		ActorUserID:  admin.ID,
		TargetUserID: target.ID,
		Status:       auth.UserStatusDisabled,
	})
	if err != nil {
		t.Fatalf("expected status update to succeed, got %v", err)
	}
	if statusResult.User.Status != auth.UserStatusDisabled {
		t.Fatalf("expected status disabled, got %s", statusResult.User.Status)
	}

	if err := svc.ResetUserPassword(context.Background(), service.ResetUserPasswordInput{
		ActorUserID:  admin.ID,
		TargetUserID: target.ID,
		NewPassword:  "ResetPass123!",
	}); err != nil {
		t.Fatalf("expected password reset to succeed, got %v", err)
	}
	credential, err := deps.Store.CredentialLocals.GetByUserID(context.Background(), target.ID)
	if err != nil {
		t.Fatalf("expected reset to create credential, got %v", err)
	}
	if err := auth.VerifyPassword(credential.PasswordHash, "ResetPass123!"); err != nil {
		t.Fatalf("expected reset password to verify, got %v", err)
	}

	audits, err := svc.ListIdentityAudits(context.Background(), service.ListIdentityAuditsInput{
		ActorUserID: admin.ID,
		EventType:   auth.AuditEventAdminUpdateUserStatus,
		Page:        1,
		PageSize:    10,
	})
	if err != nil {
		t.Fatalf("expected audit list to succeed, got %v", err)
	}
	if audits.Total != 1 || len(audits.Items) != 1 {
		t.Fatalf("expected one status audit, got total=%d items=%d", audits.Total, len(audits.Items))
	}
}

func createUserWithRole(t *testing.T, store *repo.Store, role, status string) *entity.User {
	t.Helper()
	return testkit.CreateUser(t, store, func(user *entity.User) {
		user.Role = role
		user.Status = status
	})
}

func stringPtr(value string) *string {
	return &value
}
