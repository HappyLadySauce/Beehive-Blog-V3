package auth_test

import (
	"testing"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
)

// TestIssueAndParseAccessToken verifies access token issue and parse behavior.
// TestIssueAndParseAccessToken 验证 access token 的签发与解析行为。
func TestIssueAndParseAccessToken(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 4, 20, 10, 0, 0, 0, time.UTC)
	token, expiresAt, err := auth.IssueAccessToken(
		"secret",
		15*time.Minute,
		101,
		auth.UserRoleMember,
		auth.UserStatusActive,
		301,
		auth.AuthSourceLocal,
		now,
	)
	if err != nil {
		t.Fatalf("expected issue token to succeed, got %v", err)
	}
	if token == "" {
		t.Fatalf("expected token to be non-empty")
	}
	if !expiresAt.Equal(now.Add(15 * time.Minute)) {
		t.Fatalf("expected expiresAt to match ttl, got %v", expiresAt)
	}

	claims, err := auth.ParseAccessToken("secret", token)
	if err != nil {
		t.Fatalf("expected parse token to succeed, got %v", err)
	}
	if claims.UserID != 101 || claims.SessionID != 301 {
		t.Fatalf("unexpected claims: %+v", claims)
	}
}

// TestParseAccessTokenRejectsWrongSecret verifies signature validation.
// TestParseAccessTokenRejectsWrongSecret 验证签名密钥错误时的拒绝行为。
func TestParseAccessTokenRejectsWrongSecret(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	token, _, err := auth.IssueAccessToken("secret", time.Minute, 1, auth.UserRoleMember, auth.UserStatusActive, 2, auth.AuthSourceLocal, now)
	if err != nil {
		t.Fatalf("expected issue token to succeed, got %v", err)
	}

	if _, err := auth.ParseAccessToken("wrong-secret", token); err == nil {
		t.Fatalf("expected parse with wrong secret to fail")
	}
}

// TestGenerateRefreshToken verifies refresh token generation and hashing.
// TestGenerateRefreshToken 验证 refresh token 的生成与哈希。
func TestGenerateRefreshToken(t *testing.T) {
	t.Parallel()

	token1, err := auth.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("expected generate refresh token to succeed, got %v", err)
	}
	token2, err := auth.GenerateRefreshToken()
	if err != nil {
		t.Fatalf("expected generate refresh token to succeed, got %v", err)
	}
	if token1 == "" || token2 == "" {
		t.Fatalf("expected refresh tokens to be non-empty")
	}
	if token1 == token2 {
		t.Fatalf("expected refresh tokens to be unique")
	}

	hash1 := auth.HashRefreshToken(token1)
	hash2 := auth.HashRefreshToken(token1)
	if hash1 == "" {
		t.Fatalf("expected token hash to be non-empty")
	}
	if hash1 != hash2 {
		t.Fatalf("expected token hashing to be deterministic")
	}
}
