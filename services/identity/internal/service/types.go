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

// UserListResult contains admin user list output.
// UserListResult 包含管理员用户列表输出。
type UserListResult struct {
	Items    []entity.User
	Total    int64
	Page     int
	PageSize int
}

// AdminUserResult contains a managed user snapshot.
// AdminUserResult 包含被管理用户快照。
type AdminUserResult struct {
	User *entity.User
}

// AuditListResult contains identity audit list output.
// AuditListResult 包含 identity 审计列表输出。
type AuditListResult struct {
	Items    []entity.IdentityAudit
	Total    int64
	Page     int
	PageSize int
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

// PingResult contains identity readiness information.
// PingResult 包含 identity 就绪信息。
type PingResult struct {
	OK      bool
	Service string
	Version string
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

// ListUsersInput describes admin user-list filters.
// ListUsersInput 描述管理员用户列表过滤参数。
type ListUsersInput struct {
	ActorUserID int64
	Keyword     string
	Role        string
	Status      string
	Page        int
	PageSize    int
}

// UpdateOwnProfileInput describes self-service profile updates.
// UpdateOwnProfileInput 描述用户自助资料更新。
type UpdateOwnProfileInput struct {
	UserID    int64
	Nickname  string
	AvatarURL string
	ClientIP  string
}

// ChangeOwnPasswordInput describes self-service password changes.
// ChangeOwnPasswordInput 描述用户自助密码修改。
type ChangeOwnPasswordInput struct {
	UserID      int64
	OldPassword string
	NewPassword string
	ClientIP    string
}

// UpdateUserRoleInput describes admin role changes.
// UpdateUserRoleInput 描述管理员修改用户角色。
type UpdateUserRoleInput struct {
	ActorUserID  int64
	TargetUserID int64
	Role         string
	ClientIP     string
}

// UpdateUserStatusInput describes admin status changes.
// UpdateUserStatusInput 描述管理员修改用户状态。
type UpdateUserStatusInput struct {
	ActorUserID  int64
	TargetUserID int64
	Status       string
	ClientIP     string
}

// ResetUserPasswordInput describes admin password resets.
// ResetUserPasswordInput 描述管理员重置用户密码。
type ResetUserPasswordInput struct {
	ActorUserID  int64
	TargetUserID int64
	NewPassword  string
	ClientIP     string
}

// ListIdentityAuditsInput describes audit list filters.
// ListIdentityAuditsInput 描述审计列表过滤参数。
type ListIdentityAuditsInput struct {
	ActorUserID int64
	EventType   string
	Result      string
	UserID      *int64
	StartedAt   *time.Time
	EndedAt     *time.Time
	Page        int
	PageSize    int
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
