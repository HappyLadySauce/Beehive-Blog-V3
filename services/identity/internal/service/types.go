package service

import (
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
)

// AuthResult contains the authenticated user, session, and issued token material.
// AuthResult 包含认证后的用户、会话与签发的 token 材料。
type AuthResult struct {
	User            *entity.User
	Session         *entity.UserSession
	AccessToken     string
	RefreshToken    string
	AccessExpiresAt time.Time
}

// CurrentUserResult contains the trusted current user snapshot.
// CurrentUserResult 包含可信当前用户快照。
type CurrentUserResult struct {
	User *entity.User
}

// IntrospectionResult contains access token introspection output.
// IntrospectionResult 包含 access token introspection 输出。
type IntrospectionResult struct {
	Active    bool
	User      *entity.User
	Session   *entity.UserSession
	ExpiresAt int64
}

// StartSSOResult contains the outbound authorize URL payload.
// StartSSOResult 包含对外授权地址载荷。
type StartSSOResult struct {
	Provider string
	AuthURL  string
	State    string
}

// RegisterLocalUserInput describes local registration input.
// RegisterLocalUserInput 描述本地注册输入。
type RegisterLocalUserInput struct {
	Username string
	Email    string
	Password string
	Nickname string
	ClientIP string
}

// LoginLocalUserInput describes local login input.
// LoginLocalUserInput 描述本地登录输入。
type LoginLocalUserInput struct {
	LoginIdentifier string
	Password        string
	ClientType      string
	DeviceID        string
	DeviceName      string
	UserAgent       string
	ClientIP        string
}

// RefreshSessionTokenInput describes refresh input.
// RefreshSessionTokenInput 描述 refresh 输入。
type RefreshSessionTokenInput struct {
	RefreshToken string
	UserAgent    string
	ClientIP     string
}

// LogoutSessionInput describes logout input.
// LogoutSessionInput 描述登出输入。
type LogoutSessionInput struct {
	SessionID int64
	ClientIP  string
}

// GetCurrentUserInput describes current-user lookup input.
// GetCurrentUserInput 描述当前用户查询输入。
type GetCurrentUserInput struct {
	UserID int64
}

// IntrospectAccessTokenInput describes introspection input.
// IntrospectAccessTokenInput 描述 introspection 输入。
type IntrospectAccessTokenInput struct {
	AccessToken string
}

// StartSSOInput describes SSO start input.
// StartSSOInput 描述 SSO start 输入。
type StartSSOInput struct {
	Provider    string
	RedirectURI string
	State       string
	ClientIP    string
}

// FinishSSOInput describes SSO callback completion input.
// FinishSSOInput 描述 SSO callback 完成输入。
type FinishSSOInput struct {
	Provider    string
	Code        string
	State       string
	RedirectURI string
	ClientType  string
	DeviceID    string
	DeviceName  string
	UserAgent   string
	ClientIP    string
}
