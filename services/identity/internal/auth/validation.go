package auth

import (
	"fmt"
	"net/mail"
	"regexp"
	"strings"
)

var usernamePattern = regexp.MustCompile(`^[a-zA-Z0-9_]{3,32}$`)

// NormalizeUsername trims and validates a username.
// NormalizeUsername 规范化并校验用户名。
func NormalizeUsername(username string) (string, error) {
	normalized := strings.TrimSpace(username)
	if !usernamePattern.MatchString(normalized) {
		return "", fmt.Errorf("username must be 3-32 characters and contain only letters, digits, or underscores")
	}

	return normalized, nil
}

// NormalizeEmail trims and validates an email address.
// NormalizeEmail 规范化并校验邮箱地址。
func NormalizeEmail(email string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(email))
	if normalized == "" {
		return "", nil
	}

	if _, err := mail.ParseAddress(normalized); err != nil {
		return "", fmt.Errorf("email is invalid: %w", err)
	}

	return normalized, nil
}

// NormalizeNickname trims a nickname and limits its size.
// NormalizeNickname 规范化昵称并限制长度。
func NormalizeNickname(nickname string) (string, error) {
	normalized := strings.TrimSpace(nickname)
	if len(normalized) > 128 {
		return "", fmt.Errorf("nickname must not exceed 128 characters")
	}

	return normalized, nil
}

// ValidatePassword checks the minimum password policy.
// ValidatePassword 校验密码最小安全策略。
func ValidatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}

	return nil
}

// NormalizeLoginIdentifier trims and normalizes a login identifier.
// NormalizeLoginIdentifier 规范化登录标识。
func NormalizeLoginIdentifier(identifier string) (string, error) {
	normalized := strings.TrimSpace(identifier)
	if normalized == "" {
		return "", fmt.Errorf("login identifier is required")
	}

	if strings.Contains(normalized, "@") {
		return NormalizeEmail(normalized)
	}

	return NormalizeUsername(normalized)
}

// NormalizeProvider trims and validates a provider identifier.
// NormalizeProvider 规范化并校验 provider 标识。
func NormalizeProvider(provider string) (string, error) {
	normalized := strings.ToLower(strings.TrimSpace(provider))
	switch normalized {
	case ProviderGitHub:
		return normalized, nil
	default:
		return "", fmt.Errorf("unsupported provider")
	}
}
