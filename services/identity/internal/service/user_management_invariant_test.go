package service

import (
	"errors"
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
)

func TestEnsureActiveAdminInvariantRejectsRemovingOnlyActiveAdmin(t *testing.T) {
	target := &entity.User{ID: 1, Role: auth.UserRoleAdmin, Status: auth.UserStatusActive}
	activeAdmins := []entity.User{*target}

	err := ensureActiveAdminInvariant(activeAdmins, target, auth.UserRoleMember, target.Status)
	if !errors.Is(err, errs.E(errs.CodeIdentityInvalidArgument)) {
		t.Fatalf("expected invalid argument when removing only active admin role, got %v", err)
	}

	err = ensureActiveAdminInvariant(activeAdmins, target, target.Role, auth.UserStatusDisabled)
	if !errors.Is(err, errs.E(errs.CodeIdentityInvalidArgument)) {
		t.Fatalf("expected invalid argument when disabling only active admin, got %v", err)
	}

	err = ensureActiveAdminInvariant(activeAdmins, target, target.Role, auth.UserStatusDeleted)
	if !errors.Is(err, errs.E(errs.CodeIdentityInvalidArgument)) {
		t.Fatalf("expected invalid argument when deleting only active admin, got %v", err)
	}
}

func TestEnsureActiveAdminInvariantAllowsNonAdminAndNoopChanges(t *testing.T) {
	target := &entity.User{ID: 2, Role: auth.UserRoleAdmin, Status: auth.UserStatusActive}
	activeAdmins := []entity.User{
		{ID: 1, Role: auth.UserRoleAdmin, Status: auth.UserStatusActive},
		*target,
	}

	if err := ensureActiveAdminInvariant(activeAdmins, target, auth.UserRoleMember, target.Status); err != nil {
		t.Fatalf("expected removing one of multiple active admins to succeed, got %v", err)
	}
	if err := ensureActiveAdminInvariant(activeAdmins[:1], target, auth.UserRoleAdmin, target.Status); err != nil {
		t.Fatalf("expected unchanged active admin role to succeed, got %v", err)
	}
	member := &entity.User{ID: 3, Role: auth.UserRoleMember, Status: auth.UserStatusActive}
	if err := ensureActiveAdminInvariant(activeAdmins[:1], member, member.Role, auth.UserStatusDisabled); err != nil {
		t.Fatalf("expected non-admin status update to succeed, got %v", err)
	}
}
