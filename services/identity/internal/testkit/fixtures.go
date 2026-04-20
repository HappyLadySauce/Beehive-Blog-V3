package testkit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
)

// CreateUser creates a test user with sensible defaults.
// CreateUser 使用合理默认值创建测试用户。
func CreateUser(t *testing.T, store *repo.Store, opts ...func(*entity.User)) *entity.User {
	t.Helper()

	email := fmt.Sprintf("user_%d@example.com", time.Now().UnixNano())
	user := &entity.User{
		Username: fmt.Sprintf("user_%d", time.Now().UnixNano()),
		Email:    &email,
		Role:     auth.UserRoleMember,
		Status:   auth.UserStatusActive,
	}
	for _, opt := range opts {
		opt(user)
	}
	if err := store.Users.Create(context.Background(), user); err != nil {
		t.Fatalf("failed to create user fixture: %v", err)
	}

	return user
}

// CreateCredentialLocal creates a local credential fixture.
// CreateCredentialLocal 创建本地凭证夹具。
func CreateCredentialLocal(t *testing.T, store *repo.Store, userID int64, passwordHash string) *entity.CredentialLocal {
	t.Helper()

	credential := &entity.CredentialLocal{
		UserID:            userID,
		PasswordHash:      passwordHash,
		PasswordUpdatedAt: time.Now().UTC(),
	}
	if err := store.CredentialLocals.Create(context.Background(), credential); err != nil {
		t.Fatalf("failed to create credential fixture: %v", err)
	}

	return credential
}

// CreateSession creates a session fixture.
// CreateSession 创建会话夹具。
func CreateSession(t *testing.T, store *repo.Store, userID int64, opts ...func(*entity.UserSession)) *entity.UserSession {
	t.Helper()

	now := time.Now().UTC()
	session := &entity.UserSession{
		UserID:     userID,
		AuthSource: auth.AuthSourceLocal,
		Status:     auth.SessionStatusActive,
		LastSeenAt: &now,
		ExpiresAt:  now.Add(time.Hour),
	}
	for _, opt := range opts {
		opt(session)
	}
	if err := store.UserSessions.Create(context.Background(), session); err != nil {
		t.Fatalf("failed to create session fixture: %v", err)
	}

	return session
}

// CreateRefreshToken creates a refresh token fixture.
// CreateRefreshToken 创建 refresh token 夹具。
func CreateRefreshToken(t *testing.T, store *repo.Store, sessionID int64, tokenHash string, opts ...func(*entity.RefreshToken)) *entity.RefreshToken {
	t.Helper()

	now := time.Now().UTC()
	token := &entity.RefreshToken{
		SessionID: sessionID,
		TokenHash: tokenHash,
		IssuedAt:  now,
		ExpiresAt: now.Add(time.Hour),
	}
	for _, opt := range opts {
		opt(token)
	}
	if err := store.RefreshTokens.Create(context.Background(), token); err != nil {
		t.Fatalf("failed to create refresh token fixture: %v", err)
	}

	return token
}

// CreateOAuthState creates an oauth state fixture.
// CreateOAuthState 创建 oauth state 夹具。
func CreateOAuthState(t *testing.T, store *repo.Store, provider, state, redirectURI string, opts ...func(*entity.OAuthLoginState)) *entity.OAuthLoginState {
	t.Helper()

	row := &entity.OAuthLoginState{
		Provider:    provider,
		State:       state,
		RedirectURI: redirectURI,
		ExpiresAt:   time.Now().UTC().Add(10 * time.Minute),
	}
	for _, opt := range opts {
		opt(row)
	}
	if err := store.OAuthLoginStates.Create(context.Background(), row); err != nil {
		t.Fatalf("failed to create oauth state fixture: %v", err)
	}

	return row
}

// CreateFederatedIdentity creates a federated identity fixture.
// CreateFederatedIdentity 创建联邦身份夹具。
func CreateFederatedIdentity(t *testing.T, store *repo.Store, userID int64, providerName, subject string, opts ...func(*entity.FederatedIdentity)) *entity.FederatedIdentity {
	t.Helper()

	fed := &entity.FederatedIdentity{
		UserID:              userID,
		Provider:            providerName,
		ProviderSubject:     subject,
		ProviderSubjectType: "github_user_id",
	}
	for _, opt := range opts {
		opt(fed)
	}
	if err := store.FederatedIdentities.Create(context.Background(), fed); err != nil {
		t.Fatalf("failed to create federated identity fixture: %v", err)
	}

	return fed
}
