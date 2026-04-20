package auth_test

import (
	"testing"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
)

// TestNormalizeUsername verifies username normalization and validation.
// TestNormalizeUsername 验证用户名规范化与校验。
func TestNormalizeUsername(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "valid", input: "alice_123", want: "alice_123"},
		{name: "trim spaces", input: "  alice_123  ", want: "alice_123"},
		{name: "too short", input: "ab", wantErr: true},
		{name: "invalid chars", input: "alice-123", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := auth.NormalizeUsername(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

// TestNormalizeEmail verifies email normalization behavior.
// TestNormalizeEmail 验证邮箱规范化行为。
func TestNormalizeEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "valid", input: "Alice@example.COM", want: "alice@example.com"},
		{name: "empty is allowed", input: "   ", want: ""},
		{name: "invalid", input: "not-an-email", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := auth.NormalizeEmail(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

// TestNormalizeLoginIdentifier verifies username and email login identifier handling.
// TestNormalizeLoginIdentifier 验证用户名和邮箱登录标识处理。
func TestNormalizeLoginIdentifier(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{name: "username", input: " Alice_001 ", want: "Alice_001"},
		{name: "email", input: " Alice@example.com ", want: "alice@example.com"},
		{name: "empty", input: " ", wantErr: true},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := auth.NormalizeLoginIdentifier(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("expected nil error, got %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

// TestValidatePassword verifies the minimum password policy.
// TestValidatePassword 验证最小密码策略。
func TestValidatePassword(t *testing.T) {
	t.Parallel()

	if err := auth.ValidatePassword("12345678"); err != nil {
		t.Fatalf("expected valid password, got %v", err)
	}
	if err := auth.ValidatePassword("1234567"); err == nil {
		t.Fatalf("expected short password to fail")
	}
}

// TestNormalizeProvider verifies provider normalization rules.
// TestNormalizeProvider 验证 provider 规范化规则。
func TestNormalizeProvider(t *testing.T) {
	t.Parallel()

	got, err := auth.NormalizeProvider(" GitHub ")
	if err != nil {
		t.Fatalf("expected github provider to normalize, got %v", err)
	}
	if got != auth.ProviderGitHub {
		t.Fatalf("expected %q, got %q", auth.ProviderGitHub, got)
	}

	if _, err := auth.NormalizeProvider("unknown"); err == nil {
		t.Fatalf("expected unknown provider to fail")
	}
}
