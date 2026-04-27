package service

import (
	"context"
	"net/url"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
)

const (
	maxAvatarURLLength  = 2048
	maxAuditEventLength = 64
)

// UserManagementService handles identity user administration and self-service mutations.
// UserManagementService 处理 identity 用户管理与用户自助修改。
type UserManagementService struct {
	deps Dependencies
}

// NewUserManagementService creates a UserManagementService.
// NewUserManagementService 创建 UserManagementService。
func NewUserManagementService(deps Dependencies) *UserManagementService {
	return &UserManagementService{deps: deps}
}

// ListUsers returns paged users for active administrators.
// ListUsers 为活跃管理员返回分页用户列表。
func (s *UserManagementService) ListUsers(ctx context.Context, in ListUsersInput) (*UserListResult, error) {
	if _, err := s.requireActiveAdmin(ctx, in.ActorUserID); err != nil {
		return nil, err
	}

	role, err := normalizeOptionalRole(in.Role)
	if err != nil {
		return nil, err
	}
	status, err := normalizeOptionalStatus(in.Status, true)
	if err != nil {
		return nil, err
	}

	page, pageSize := normalizePageInput(in.Page, in.PageSize)
	users, total, err := s.deps.Store.Users.List(ctx, repo.ListFilter{
		Keyword:  in.Keyword,
		Role:     role,
		Status:   status,
		Page:     page,
		PageSize: pageSize,
	})
	if err != nil {
		return nil, errs.Wrap(err, errs.CodeIdentityInternal, "list users failed")
	}

	return &UserListResult{Items: users, Total: total, Page: page, PageSize: pageSize}, nil
}

// UpdateOwnProfile updates nickname and avatar URL for the current user.
// UpdateOwnProfile 更新当前用户昵称与头像地址。
func (s *UserManagementService) UpdateOwnProfile(ctx context.Context, in UpdateOwnProfileInput) (*CurrentUserResult, error) {
	var patch repo.ProfileUpdate
	auditDetail := map[string]any{}
	if in.Nickname != nil {
		nickname, err := auth.NormalizeNickname(*in.Nickname)
		if err != nil {
			return nil, errs.Wrap(err, errs.CodeIdentityInvalidArgument, "nickname is invalid")
		}
		patch.NicknameSet = true
		patch.Nickname = optionalString(nickname)
		auditDetail["nickname"] = nickname
	}
	if in.AvatarURL != nil {
		avatarURL, err := normalizeAvatarURL(*in.AvatarURL)
		if err != nil {
			return nil, err
		}
		patch.AvatarURLSet = true
		patch.AvatarURL = optionalString(avatarURL)
		auditDetail["avatar_url"] = avatarURL
	}

	now := s.deps.Clock()
	var updated *entity.User
	if err := withTransaction(ctx, s.deps.Store, func(txStore *repo.Store) error {
		user, loadErr := txStore.Users.GetByID(ctx, in.UserID)
		if loadErr != nil {
			if repo.IsNotFound(loadErr) {
				return errs.New(errs.CodeIdentityUserNotFound, "user not found")
			}
			return errs.Wrap(loadErr, errs.CodeIdentityInternal, "load user failed")
		}
		if err := validateActiveUserStatus(user.Status); err != nil {
			return err
		}

		updatedUser, updateErr := txStore.Users.UpdateProfile(ctx, in.UserID, patch, now)
		if updateErr != nil {
			return errs.Wrap(updateErr, errs.CodeIdentityInternal, "update profile failed")
		}
		updated = updatedUser
		writeAudit(ctx, txStore, auditInput{
			UserID:    &in.UserID,
			EventType: auth.AuditEventUpdateOwnProfile,
			Result:    auth.AuditResultSuccess,
			ClientIP:  stringPtr(in.ClientIP),
			Detail:    auth.MarshalAuditDetail(auditDetail),
		})
		return nil
	}); err != nil {
		return nil, err
	}

	return &CurrentUserResult{User: updated}, nil
}

// ChangeOwnPassword changes the current user's local password after old-password verification.
// ChangeOwnPassword 在校验旧密码后修改当前用户本地密码。
func (s *UserManagementService) ChangeOwnPassword(ctx context.Context, in ChangeOwnPasswordInput) error {
	if strings.TrimSpace(in.OldPassword) == "" {
		return errs.New(errs.CodeIdentityInvalidArgument, "old_password is required")
	}
	if err := auth.ValidatePassword(in.NewPassword); err != nil {
		return errs.Wrap(err, errs.CodeIdentityInvalidArgument, "new_password is invalid")
	}

	now := s.deps.Clock()
	return withTransaction(ctx, s.deps.Store, func(txStore *repo.Store) error {
		user, err := txStore.Users.GetByID(ctx, in.UserID)
		if err != nil {
			if repo.IsNotFound(err) {
				return errs.New(errs.CodeIdentityUserNotFound, "user not found")
			}
			return errs.Wrap(err, errs.CodeIdentityInternal, "load user failed")
		}
		if err := validateActiveUserStatus(user.Status); err != nil {
			return err
		}
		credential, err := txStore.CredentialLocals.GetByUserID(ctx, in.UserID)
		if err != nil {
			if repo.IsNotFound(err) {
				return errs.New(errs.CodeIdentityInvalidCredentials, "local credential not found")
			}
			return errs.Wrap(err, errs.CodeIdentityInternal, "load credential failed")
		}
		if err := auth.VerifyPassword(credential.PasswordHash, in.OldPassword); err != nil {
			return errs.Wrap(err, errs.CodeIdentityInvalidCredentials, "invalid credentials")
		}
		hash, err := auth.HashPassword(in.NewPassword, s.deps.Config.Security.PasswordHashCost)
		if err != nil {
			return errs.Wrap(err, errs.CodeIdentityInternal, "hash password failed")
		}
		if err := txStore.CredentialLocals.UpdatePasswordHash(ctx, in.UserID, hash, now); err != nil {
			return errs.Wrap(err, errs.CodeIdentityInternal, "update password failed")
		}
		writeAudit(ctx, txStore, auditInput{
			UserID:    &in.UserID,
			EventType: auth.AuditEventChangeOwnPassword,
			Result:    auth.AuditResultSuccess,
			ClientIP:  stringPtr(in.ClientIP),
		})
		return nil
	})
}

// UpdateUserRole updates a target user's role by an active administrator.
// UpdateUserRole 由活跃管理员修改目标用户角色。
func (s *UserManagementService) UpdateUserRole(ctx context.Context, in UpdateUserRoleInput) (*AdminUserResult, error) {
	role, err := normalizeRequiredRole(in.Role)
	if err != nil {
		return nil, err
	}
	return s.updateManagedUser(ctx, in.ActorUserID, in.TargetUserID, in.ClientIP, auth.AuditEventAdminUpdateUserRole, func(txStore *repo.Store, target *entity.User) (*entity.User, error) {
		return txStore.Users.UpdateRole(ctx, target.ID, role, s.deps.Clock())
	}, map[string]any{"new_role": role})
}

// UpdateUserStatus updates a target user's status by an active administrator.
// UpdateUserStatus 由活跃管理员修改目标用户状态。
func (s *UserManagementService) UpdateUserStatus(ctx context.Context, in UpdateUserStatusInput) (*AdminUserResult, error) {
	status, err := normalizeOptionalStatus(in.Status, false)
	if err != nil {
		return nil, err
	}
	return s.updateManagedUser(ctx, in.ActorUserID, in.TargetUserID, in.ClientIP, auth.AuditEventAdminUpdateUserStatus, func(txStore *repo.Store, target *entity.User) (*entity.User, error) {
		return txStore.Users.UpdateStatus(ctx, target.ID, status, s.deps.Clock())
	}, map[string]any{"new_status": status})
}

// ResetUserPassword replaces a target user's local password by an active administrator.
// ResetUserPassword 由活跃管理员替换目标用户本地密码。
func (s *UserManagementService) ResetUserPassword(ctx context.Context, in ResetUserPasswordInput) error {
	if err := auth.ValidatePassword(in.NewPassword); err != nil {
		return errs.Wrap(err, errs.CodeIdentityInvalidArgument, "new_password is invalid")
	}

	_, err := s.updateManagedUser(ctx, in.ActorUserID, in.TargetUserID, in.ClientIP, auth.AuditEventAdminResetUserPassword, func(txStore *repo.Store, target *entity.User) (*entity.User, error) {
		hash, err := auth.HashPassword(in.NewPassword, s.deps.Config.Security.PasswordHashCost)
		if err != nil {
			return nil, errs.Wrap(err, errs.CodeIdentityInternal, "hash password failed")
		}
		if err := txStore.CredentialLocals.UpdatePasswordHash(ctx, target.ID, hash, s.deps.Clock()); err != nil {
			if repo.IsNotFound(err) {
				if err := txStore.CredentialLocals.Create(ctx, &entity.CredentialLocal{
					UserID:            target.ID,
					PasswordHash:      hash,
					PasswordUpdatedAt: s.deps.Clock(),
				}); err != nil {
					return nil, errs.Wrap(err, errs.CodeIdentityInternal, "create password failed")
				}
				return target, nil
			}
			return nil, errs.Wrap(err, errs.CodeIdentityInternal, "update password failed")
		}
		return target, nil
	}, nil)
	return err
}

// ListIdentityAudits returns paged audits for active administrators.
// ListIdentityAudits 为活跃管理员返回分页审计记录。
func (s *UserManagementService) ListIdentityAudits(ctx context.Context, in ListIdentityAuditsInput) (*AuditListResult, error) {
	if _, err := s.requireActiveAdmin(ctx, in.ActorUserID); err != nil {
		return nil, err
	}
	if result := strings.TrimSpace(in.Result); result != "" && result != auth.AuditResultSuccess && result != auth.AuditResultFailure {
		return nil, errs.New(errs.CodeIdentityInvalidArgument, "audit result is invalid")
	}
	if len(strings.TrimSpace(in.EventType)) > maxAuditEventLength {
		return nil, errs.New(errs.CodeIdentityInvalidArgument, "event_type must not exceed 64 characters")
	}
	if in.StartedAt != nil && in.EndedAt != nil && in.StartedAt.After(*in.EndedAt) {
		return nil, errs.New(errs.CodeIdentityInvalidArgument, "started_at must be before ended_at")
	}

	page, pageSize := normalizePageInput(in.Page, in.PageSize)
	audits, total, err := s.deps.Store.IdentityAudits.List(ctx, repo.AuditListFilter{
		EventType: strings.TrimSpace(in.EventType),
		Result:    strings.TrimSpace(in.Result),
		UserID:    in.UserID,
		StartedAt: in.StartedAt,
		EndedAt:   in.EndedAt,
		Page:      page,
		PageSize:  pageSize,
	})
	if err != nil {
		return nil, errs.Wrap(err, errs.CodeIdentityInternal, "list identity audits failed")
	}

	return &AuditListResult{Items: audits, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *UserManagementService) updateManagedUser(
	ctx context.Context,
	actorUserID int64,
	targetUserID int64,
	clientIP string,
	eventType string,
	mutate func(txStore *repo.Store, target *entity.User) (*entity.User, error),
	detail map[string]any,
) (*AdminUserResult, error) {
	if actorUserID == targetUserID {
		return nil, errs.New(errs.CodeIdentityInvalidArgument, "cannot modify own role, status, or password through admin endpoints")
	}

	var updated *entity.User
	if err := withTransaction(ctx, s.deps.Store, func(txStore *repo.Store) error {
		if _, err := s.requireActiveAdminWithStore(ctx, txStore, actorUserID); err != nil {
			return err
		}
		target, err := txStore.Users.GetByID(ctx, targetUserID)
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
	switch status {
	case auth.UserStatusActive, auth.UserStatusDisabled, auth.UserStatusLocked:
		return status, nil
	case auth.UserStatusPending:
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

func optionalString(value string) *string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func normalizeAvatarURL(value string) (string, error) {
	avatarURL := strings.TrimSpace(value)
	if avatarURL == "" {
		return "", nil
	}
	if len(avatarURL) > maxAvatarURLLength {
		return "", errs.New(errs.CodeIdentityInvalidArgument, "avatar_url must not exceed 2048 characters")
	}

	parsed, err := url.ParseRequestURI(avatarURL)
	if err != nil || parsed.Scheme == "" || parsed.Host == "" {
		return "", errs.New(errs.CodeIdentityInvalidArgument, "avatar_url is invalid")
	}
	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return "", errs.New(errs.CodeIdentityInvalidArgument, "avatar_url scheme is invalid")
	}

	return avatarURL, nil
}
