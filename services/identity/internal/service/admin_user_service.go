package service

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
)

// UpdateUserRole updates a target user's role by an active administrator.
// UpdateUserRole 由活跃管理员修改目标用户角色。
func (s *UserManagementService) UpdateUserRole(ctx context.Context, in UpdateUserRoleInput) (*AdminUserResult, error) {
	role, err := normalizeRequiredRole(in.Role)
	if err != nil {
		return nil, err
	}
	return s.updateManagedUser(ctx, in.ActorUserID, in.TargetUserID, in.ClientIP, auth.AuditEventAdminUpdateUserRole, func(txStore *repo.Store, target *entity.User) (*entity.User, error) {
		return txStore.Users.UpdateRole(ctx, target.ID, role, s.deps.Clock())
	}, func(activeAdmins []entity.User, target *entity.User) error {
		return ensureActiveAdminInvariant(activeAdmins, target, role, target.Status)
	}, map[string]any{"new_role": role})
}

// UpdateUserStatus updates a target user's status by an active administrator.
// UpdateUserStatus 由活跃管理员修改目标用户状态。
func (s *UserManagementService) UpdateUserStatus(ctx context.Context, in UpdateUserStatusInput) (*AdminUserResult, error) {
	status, err := normalizeRequiredStatus(in.Status)
	if err != nil {
		return nil, err
	}
	return s.updateManagedUser(ctx, in.ActorUserID, in.TargetUserID, in.ClientIP, auth.AuditEventAdminUpdateUserStatus, func(txStore *repo.Store, target *entity.User) (*entity.User, error) {
		return txStore.Users.UpdateStatus(ctx, target.ID, status, s.deps.Clock())
	}, func(activeAdmins []entity.User, target *entity.User) error {
		return ensureActiveAdminInvariant(activeAdmins, target, target.Role, status)
	}, map[string]any{"new_status": status})
}

func (s *UserManagementService) updateManagedUser(
	ctx context.Context,
	actorUserID int64,
	targetUserID int64,
	clientIP string,
	eventType string,
	mutate func(txStore *repo.Store, target *entity.User) (*entity.User, error),
	validate func(activeAdmins []entity.User, target *entity.User) error,
	detail map[string]any,
) (*AdminUserResult, error) {
	if actorUserID == targetUserID {
		return nil, errs.New(errs.CodeIdentityInvalidArgument, "cannot modify own role, status, or password through admin endpoints")
	}

	var updated *entity.User
	if err := withTransaction(ctx, s.deps.Store, func(txStore *repo.Store) error {
		activeAdmins, err := s.requireActiveAdminSetForUpdateWithStore(ctx, txStore, actorUserID)
		if err != nil {
			return err
		}
		target, err := txStore.Users.GetForUpdateByID(ctx, targetUserID)
		if err != nil {
			if repo.IsNotFound(err) {
				return errs.New(errs.CodeIdentityUserNotFound, "target user not found")
			}
			return errs.Wrap(err, errs.CodeIdentityInternal, "load target user failed")
		}

		if detail == nil {
			detail = map[string]any{}
		}
		detail["target_user_id"] = targetUserID
		detail["old_role"] = target.Role
		detail["old_status"] = target.Status

		if validate != nil {
			if err := validate(activeAdmins, target); err != nil {
				return err
			}
		}
		updatedUser, err := mutate(txStore, target)
		if err != nil {
			return err
		}
		updated = updatedUser
		detail["updated_role"] = updatedUser.Role
		detail["updated_status"] = updatedUser.Status
		writeAudit(ctx, txStore, auditInput{
			UserID:    &actorUserID,
			EventType: eventType,
			Result:    auth.AuditResultSuccess,
			ClientIP:  stringPtr(clientIP),
			Detail:    auth.MarshalAuditDetail(detail),
		})
		return nil
	}); err != nil {
		return nil, err
	}

	return &AdminUserResult{User: updated}, nil
}
