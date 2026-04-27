package service

import (
	"context"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
)

func (s *UserManagementService) requireActiveAdmin(ctx context.Context, userID int64) (*entity.User, error) {
	return s.requireActiveAdminWithStore(ctx, s.deps.Store, userID)
}

func (s *UserManagementService) requireActiveAdminWithStore(ctx context.Context, store *repo.Store, userID int64) (*entity.User, error) {
	if userID <= 0 {
		return nil, errs.New(errs.CodeIdentityInvalidArgument, "actor_user_id is invalid")
	}

	user, err := store.Users.GetByID(ctx, userID)
	if err != nil {
		if repo.IsNotFound(err) {
			return nil, errs.New(errs.CodeIdentityUserNotFound, "actor user not found")
		}
		return nil, errs.Wrap(err, errs.CodeIdentityInternal, "load actor user failed")
	}
	if err := validateActiveUserStatus(user.Status); err != nil {
		return nil, err
	}
	if user.Role != auth.UserRoleAdmin {
		return nil, errs.New(errs.CodeIdentityAccessForbidden, "identity administration requires admin role")
	}
	return user, nil
}

func (s *UserManagementService) requireActiveAdminSetForUpdateWithStore(ctx context.Context, store *repo.Store, userID int64) ([]entity.User, error) {
	if userID <= 0 {
		return nil, errs.New(errs.CodeIdentityInvalidArgument, "actor_user_id is invalid")
	}

	activeAdmins, err := store.Users.ListActiveAdminsForUpdate(ctx)
	if err != nil {
		return nil, errs.Wrap(err, errs.CodeIdentityInternal, "lock active admins failed")
	}

	for i := range activeAdmins {
		if activeAdmins[i].ID == userID {
			return activeAdmins, nil
		}
	}

	user, err := store.Users.GetForUpdateByID(ctx, userID)
	if err != nil {
		if repo.IsNotFound(err) {
			return nil, errs.New(errs.CodeIdentityUserNotFound, "actor user not found")
		}
		return nil, errs.Wrap(err, errs.CodeIdentityInternal, "load actor user failed")
	}
	if err := validateActiveUserStatus(user.Status); err != nil {
		return nil, err
	}
	if user.Role != auth.UserRoleAdmin {
		return nil, errs.New(errs.CodeIdentityAccessForbidden, "identity administration requires admin role")
	}
	return activeAdmins, nil
}

func ensureActiveAdminInvariant(activeAdmins []entity.User, target *entity.User, nextRole string, nextStatus string) error {
	if target == nil || target.Role != auth.UserRoleAdmin || target.Status != auth.UserStatusActive {
		return nil
	}

	nextRole = strings.ToLower(strings.TrimSpace(nextRole))
	nextStatus = strings.ToLower(strings.TrimSpace(nextStatus))
	if nextRole == auth.UserRoleAdmin && nextStatus == auth.UserStatusActive {
		return nil
	}
	if len(activeAdmins) <= 1 {
		return errs.New(errs.CodeIdentityInvalidArgument, "at least one active admin is required")
	}
	return nil
}

func normalizeRequiredRole(value string) (string, error) {
	role := strings.ToLower(strings.TrimSpace(value))
	switch role {
	case auth.UserRoleMember, auth.UserRoleAdmin:
		return role, nil
	default:
		return "", errs.New(errs.CodeIdentityInvalidArgument, "role is invalid")
	}
}

func normalizeOptionalRole(value string) (string, error) {
	role := strings.ToLower(strings.TrimSpace(value))
	if role == "" {
		return "", nil
	}
	return normalizeRequiredRole(role)
}

func normalizeOptionalStatus(value string, allowPending bool) (string, error) {
	status := strings.ToLower(strings.TrimSpace(value))
	if status == "" {
		return "", nil
	}
	return normalizeStatusValue(status, allowPending)
}

func normalizeRequiredStatus(value string) (string, error) {
	status := strings.ToLower(strings.TrimSpace(value))
	if status == "" {
		return "", errs.New(errs.CodeIdentityInvalidArgument, "status is invalid")
	}
	return normalizeStatusValue(status, false)
}

func normalizeStatusValue(status string, allowPending bool) (string, error) {
	switch status {
	case auth.UserStatusActive, auth.UserStatusDisabled, auth.UserStatusLocked:
		return status, nil
	case auth.UserStatusPending:
		if allowPending {
			return status, nil
		}
	case auth.UserStatusDeleted:
		if allowPending {
			return status, nil
		}
	}
	return "", errs.New(errs.CodeIdentityInvalidArgument, "status is invalid")
}

func normalizePageInput(page, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}
