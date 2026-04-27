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
	if !patch.NicknameSet && !patch.AvatarURLSet {
		return nil, errs.New(errs.CodeIdentityInvalidArgument, "profile patch must include at least one field")
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
