package service

import (
	"context"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
)

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
		user, err := txStore.Users.GetForUpdateByID(ctx, in.UserID)
		if err != nil {
			if repo.IsNotFound(err) {
				return errs.New(errs.CodeIdentityUserNotFound, "user not found")
			}
			return errs.Wrap(err, errs.CodeIdentityInternal, "load user failed")
		}
		if err := validateActiveUserStatus(user.Status); err != nil {
			return err
		}
		credential, err := txStore.CredentialLocals.GetForUpdateByUserID(ctx, in.UserID)
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
		if err := revokeUserSessionsAndRefreshTokens(ctx, txStore, in.UserID, now); err != nil {
			return err
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

// ResetUserPassword replaces a target user's local password by an active administrator.
// ResetUserPassword 由活跃管理员替换目标用户本地密码。
func (s *UserManagementService) ResetUserPassword(ctx context.Context, in ResetUserPasswordInput) error {
	if err := auth.ValidatePassword(in.NewPassword); err != nil {
		return errs.Wrap(err, errs.CodeIdentityInvalidArgument, "new_password is invalid")
	}

	now := s.deps.Clock()
	_, err := s.updateManagedUser(ctx, in.ActorUserID, in.TargetUserID, in.ClientIP, auth.AuditEventAdminResetUserPassword, func(txStore *repo.Store, target *entity.User) (*entity.User, error) {
		hash, err := auth.HashPassword(in.NewPassword, s.deps.Config.Security.PasswordHashCost)
		if err != nil {
			return nil, errs.Wrap(err, errs.CodeIdentityInternal, "hash password failed")
		}
		if _, err := txStore.CredentialLocals.GetForUpdateByUserID(ctx, target.ID); err != nil {
			if repo.IsNotFound(err) {
				if err := txStore.CredentialLocals.Create(ctx, &entity.CredentialLocal{
					UserID:            target.ID,
					PasswordHash:      hash,
					PasswordUpdatedAt: now,
				}); err != nil {
					return nil, errs.Wrap(err, errs.CodeIdentityInternal, "create password failed")
				}
			} else {
				return nil, errs.Wrap(err, errs.CodeIdentityInternal, "load credential failed")
			}
		} else if err := txStore.CredentialLocals.UpdatePasswordHash(ctx, target.ID, hash, now); err != nil {
			return nil, errs.Wrap(err, errs.CodeIdentityInternal, "update password failed")
		}
		if err := revokeUserSessionsAndRefreshTokens(ctx, txStore, target.ID, now); err != nil {
			return nil, err
		}
		return target, nil
	}, nil, nil)
	return err
}

func revokeUserSessionsAndRefreshTokens(ctx context.Context, txStore *repo.Store, userID int64, revokedAt time.Time) error {
	if err := txStore.RefreshTokens.RevokeActiveByUserID(ctx, userID, revokedAt); err != nil {
		return errs.Wrap(err, errs.CodeIdentityInternal, "revoke refresh tokens failed")
	}
	if err := txStore.UserSessions.RevokeActiveByUserID(ctx, userID, revokedAt); err != nil {
		return errs.Wrap(err, errs.CodeIdentityInternal, "revoke sessions failed")
	}
	return nil
}
