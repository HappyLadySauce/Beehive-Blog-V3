package service

import (
	"context"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
)

const (
	emailVerificationPassword = "password"
	emailVerificationSSO      = "sso"
)

// UpdateOwnEmail updates the current user's email after password or SSO reauthentication.
// UpdateOwnEmail 在密码或 SSO 重验通过后更新当前用户邮箱。
func (s *UserManagementService) UpdateOwnEmail(ctx context.Context, in UpdateOwnEmailInput) (*CurrentUserResult, error) {
	email, err := auth.NormalizeEmail(in.Email)
	if err != nil {
		return nil, errs.Wrap(err, errs.CodeIdentityInvalidArgument, "email is invalid")
	}
	if email == "" {
		return nil, errs.New(errs.CodeIdentityInvalidArgument, "email is required")
	}

	switch strings.ToLower(strings.TrimSpace(in.VerificationMethod)) {
	case emailVerificationPassword:
		return s.updateOwnEmailWithPassword(ctx, in, email)
	case emailVerificationSSO:
		return s.updateOwnEmailWithSSO(ctx, in, email)
	default:
		return nil, errs.New(errs.CodeIdentityInvalidArgument, "verification_method is invalid")
	}
}

func (s *UserManagementService) updateOwnEmailWithPassword(ctx context.Context, in UpdateOwnEmailInput, email string) (*CurrentUserResult, error) {
	if strings.TrimSpace(in.CurrentPassword) == "" {
		return nil, errs.New(errs.CodeIdentityInvalidArgument, "current_password is required")
	}

	now := s.deps.Clock()
	var updated *entity.User
	if err := withTransaction(ctx, s.deps.Store, func(txStore *repo.Store) error {
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
		if err := auth.VerifyPassword(credential.PasswordHash, in.CurrentPassword); err != nil {
			return errs.Wrap(err, errs.CodeIdentityInvalidCredentials, "invalid credentials")
		}
		next, err := txStore.Users.UpdateUserProfile(ctx, in.UserID, repo.UserProfileUpdate{
			EmailSet: true,
			Email:    optionalString(email),
		}, now)
		if err != nil {
			return mapEmailUpdateError(err)
		}
		updated = next
		writeAudit(ctx, txStore, auditInput{
			UserID:    &in.UserID,
			EventType: auth.AuditEventUpdateOwnEmail,
			Result:    auth.AuditResultSuccess,
			ClientIP:  stringPtr(in.ClientIP),
			Detail: auth.MarshalAuditDetail(map[string]any{
				"verification_method": emailVerificationPassword,
			}),
		})
		return nil
	}); err != nil {
		return nil, err
	}

	return &CurrentUserResult{User: updated}, nil
}

func (s *UserManagementService) updateOwnEmailWithSSO(ctx context.Context, in UpdateOwnEmailInput, email string) (*CurrentUserResult, error) {
	providerName, err := auth.NormalizeProvider(in.Provider)
	if err != nil {
		return nil, errs.Wrap(err, errs.CodeIdentityInvalidArgument, "unsupported provider")
	}
	callbackProvider, ok := s.deps.Providers.GetCallback(providerName)
	if !ok || !callbackProvider.LoginReady() {
		return nil, errs.New(errs.CodeIdentitySSOProviderNotReady, "sso provider is not ready")
	}
	if !callbackProvider.Enabled() {
		return nil, errs.New(errs.CodeIdentitySSOProviderDisabled, "sso provider is disabled")
	}
	redirectURI := strings.TrimSpace(in.RedirectURI)
	if strings.TrimSpace(in.Code) == "" || strings.TrimSpace(in.State) == "" || redirectURI == "" {
		return nil, errs.New(errs.CodeIdentityInvalidArgument, "code, state, and redirect_uri are required")
	}
	if redirectURI != strings.TrimSpace(callbackProvider.RedirectURL()) {
		return nil, errs.New(errs.CodeIdentityInvalidArgument, "redirect_uri does not match configured provider redirect")
	}

	stateID, err := s.validateReauthState(ctx, providerName, strings.TrimSpace(in.State), redirectURI, in.UserID)
	if err != nil {
		return nil, err
	}
	accessToken, err := callbackProvider.ExchangeCode(ctx, in.Code, redirectURI)
	if err != nil {
		return nil, errs.Wrap(err, errs.CodeIdentityInvalidCredentials, "exchange provider code failed")
	}
	profile, _, err := callbackProvider.FetchProfile(ctx, accessToken)
	if err != nil {
		return nil, errs.Wrap(err, errs.CodeIdentityInvalidCredentials, "fetch provider profile failed")
	}

	now := s.deps.Clock()
	var updated *entity.User
	if err := withTransaction(ctx, s.deps.Store, func(txStore *repo.Store) error {
		lockedState, err := txStore.OAuthLoginStates.GetForUpdateByProviderState(ctx, providerName, strings.TrimSpace(in.State))
		if err != nil {
			if repo.IsNotFound(err) {
				return newSSOStateFailure("state_not_found", errs.CodeIdentitySSOStateInvalid, "sso state is invalid")
			}
			return err
		}
		if lockedState.ID != stateID || lockedState.ConsumedAt != nil || lockedState.ExpiresAt.Before(now) {
			return newSSOStateFailure("state_invalid", errs.CodeIdentitySSOStateInvalid, "sso state is invalid")
		}
		if lockedState.RedirectURI != redirectURI || lockedState.Purpose != oauthStatePurposeEmailUpdate || lockedState.SubjectUserID == nil || *lockedState.SubjectUserID != in.UserID {
			return newSSOStateFailure("state_scope_mismatch", errs.CodeIdentitySSOStateInvalid, "sso state is invalid")
		}
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
		fed, err := txStore.FederatedIdentities.GetByProviderIdentity(ctx, providerName, profile.Subject, profile.OpenID, profile.UnionID)
		if err != nil {
			if repo.IsNotFound(err) {
				return errs.New(errs.CodeIdentityInvalidCredentials, "sso identity is not linked to current user")
			}
			return err
		}
		if fed.UserID != in.UserID {
			return errs.New(errs.CodeIdentityInvalidCredentials, "sso identity does not match current user")
		}
		next, err := txStore.Users.UpdateUserProfile(ctx, in.UserID, repo.UserProfileUpdate{
			EmailSet: true,
			Email:    optionalString(email),
		}, now)
		if err != nil {
			return mapEmailUpdateError(err)
		}
		if err := txStore.OAuthLoginStates.Consume(ctx, lockedState.ID, now); err != nil {
			return err
		}
		updated = next
		writeAudit(ctx, txStore, auditInput{
			UserID:    &in.UserID,
			Provider:  stringPtr(providerName),
			EventType: auth.AuditEventUpdateOwnEmail,
			Result:    auth.AuditResultSuccess,
			ClientIP:  stringPtr(in.ClientIP),
			Detail: auth.MarshalAuditDetail(map[string]any{
				"verification_method": emailVerificationSSO,
			}),
		})
		return nil
	}); err != nil {
		return nil, err
	}

	return &CurrentUserResult{User: updated}, nil
}

func (s *UserManagementService) validateReauthState(ctx context.Context, providerName, state, redirectURI string, userID int64) (int64, error) {
	now := s.deps.Clock()
	var stateID int64
	if err := withTransaction(ctx, s.deps.Store, func(txStore *repo.Store) error {
		stateRow, err := txStore.OAuthLoginStates.GetForUpdateByProviderState(ctx, providerName, state)
		if err != nil {
			if repo.IsNotFound(err) {
				return newSSOStateFailure("state_not_found", errs.CodeIdentitySSOStateInvalid, "sso state is invalid")
			}
			return err
		}
		if stateRow.ConsumedAt != nil || stateRow.ExpiresAt.Before(now) {
			return newSSOStateFailure("state_invalid", errs.CodeIdentitySSOStateInvalid, "sso state is invalid")
		}
		if stateRow.RedirectURI != redirectURI || stateRow.Purpose != oauthStatePurposeEmailUpdate || stateRow.SubjectUserID == nil || *stateRow.SubjectUserID != userID {
			return newSSOStateFailure("state_scope_mismatch", errs.CodeIdentitySSOStateInvalid, "sso state is invalid")
		}
		stateID = stateRow.ID
		return nil
	}); err != nil {
		return 0, err
	}

	return stateID, nil
}

func mapEmailUpdateError(err error) error {
	if conflictKind, ok := repo.ParseUniqueViolation(err); ok {
		if conflictKind == repo.UniqueViolationEmail {
			return errs.Wrap(err, errs.CodeIdentityEmailAlreadyExists, "email already exists")
		}
	}
	return errs.Wrap(err, errs.CodeIdentityInternal, "update email failed")
}
