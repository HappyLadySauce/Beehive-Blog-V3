package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// AccessTokenClaims defines the JWT claims used by identity.
// AccessTokenClaims 定义 identity 使用的 JWT claims。
type AccessTokenClaims struct {
	UserID        int64  `json:"user_id"`
	Role          string `json:"role"`
	AccountStatus string `json:"account_status"`
	SessionID     int64  `json:"session_id"`
	AuthSource    string `json:"auth_source"`
	jwt.RegisteredClaims
}

// IssueAccessToken signs a short-lived access token.
// IssueAccessToken 签发短期 access token。
func IssueAccessToken(secret string, ttl time.Duration, userID int64, role, accountStatus string, sessionID int64, authSource string, now time.Time) (string, time.Time, error) {
	expiresAt := now.Add(ttl)
	claims := AccessTokenClaims{
		UserID:        userID,
		Role:          role,
		AccountStatus: accountStatus,
		SessionID:     sessionID,
		AuthSource:    authSource,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(userID, 10),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, expiresAt, nil
}

// ParseAccessToken parses an access token without relying on built-in exp validation.
// ParseAccessToken 解析 access token，并手动处理过期时间校验。
func ParseAccessToken(secret, accessToken string) (*AccessTokenClaims, error) {
	claims := &AccessTokenClaims{}
	parser := jwt.Parser{
		SkipClaimsValidation: true,
		ValidMethods:         []string{jwt.SigningMethodHS256.Alg()},
	}

	token, err := parser.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, fmt.Errorf("access token is invalid")
	}

	return claims, nil
}

// GenerateRefreshToken creates a high-entropy opaque refresh token.
// GenerateRefreshToken 生成高熵 opaque refresh token。
func GenerateRefreshToken() (string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// HashRefreshToken hashes a refresh token for database storage.
// HashRefreshToken 对 refresh token 做哈希以便数据库存储。
func HashRefreshToken(refreshToken string) string {
	sum := sha256.Sum256([]byte(refreshToken))
	return hex.EncodeToString(sum[:])
}
