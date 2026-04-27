package auth

const (
	// ProviderGitHub identifies the GitHub SSO provider.
	// ProviderGitHub 标识 GitHub SSO provider。
	ProviderGitHub = "github"
	// ProviderQQ identifies the QQ SSO provider.
	// ProviderQQ 标识 QQ SSO provider。
	ProviderQQ = "qq"
	// ProviderWeChat identifies the WeChat SSO provider.
	// ProviderWeChat 标识 WeChat SSO provider。
	ProviderWeChat = "wechat"
)

const (
	// AuthSourceLocal marks local authentication sessions.
	// AuthSourceLocal 标记本地认证会话。
	AuthSourceLocal = "local"
	// AuthSourceSSO marks SSO authentication sessions.
	// AuthSourceSSO 标记 SSO 认证会话。
	AuthSourceSSO = "sso"
)

const (
	// UserRoleMember is the default role for public registrations.
	// UserRoleMember 是公开注册用户的默认角色。
	UserRoleMember = "member"
	// UserRoleAdmin is the platform administrator role.
	// UserRoleAdmin 是平台管理员角色。
	UserRoleAdmin = "admin"
)

const (
	// UserStatusPending means the account is pending activation.
	// UserStatusPending 表示账号待激活。
	UserStatusPending = "pending"
	// UserStatusActive means the account can authenticate.
	// UserStatusActive 表示账号可正常认证。
	UserStatusActive = "active"
	// UserStatusDisabled means the account is disabled.
	// UserStatusDisabled 表示账号已禁用。
	UserStatusDisabled = "disabled"
	// UserStatusLocked means the account is locked.
	// UserStatusLocked 表示账号已锁定。
	UserStatusLocked = "locked"
)

const (
	// SessionStatusActive means the session is active.
	// SessionStatusActive 表示会话处于活跃状态。
	SessionStatusActive = "active"
	// SessionStatusRevoked means the session is revoked.
	// SessionStatusRevoked 表示会话已吊销。
	SessionStatusRevoked = "revoked"
	// SessionStatusExpired means the session is expired.
	// SessionStatusExpired 表示会话已过期。
	SessionStatusExpired = "expired"
)

const (
	// AuditResultSuccess marks successful audit outcomes.
	// AuditResultSuccess 标记成功审计结果。
	AuditResultSuccess = "success"
	// AuditResultFailure marks failed audit outcomes.
	// AuditResultFailure 标记失败审计结果。
	AuditResultFailure = "failure"
)

const (
	// AuditEventRegisterLocal marks local registration.
	// AuditEventRegisterLocal 标记本地注册事件。
	AuditEventRegisterLocal = "register_local_user"
	// AuditEventLoginLocal marks local login.
	// AuditEventLoginLocal 标记本地登录事件。
	AuditEventLoginLocal = "login_local_user"
	// AuditEventRefreshSession marks refresh rotation.
	// AuditEventRefreshSession 标记 refresh 轮换事件。
	AuditEventRefreshSession = "refresh_session_token"
	// AuditEventLogoutSession marks session logout.
	// AuditEventLogoutSession 标记会话登出事件。
	AuditEventLogoutSession = "logout_session"
	// AuditEventStartSSO marks SSO start.
	// AuditEventStartSSO 标记 SSO 开始事件。
	AuditEventStartSSO = "start_sso_login"
	// AuditEventFinishSSO marks SSO callback completion.
	// AuditEventFinishSSO 标记 SSO 回调完成事件。
	AuditEventFinishSSO = "finish_sso_login"
	// AuditEventAdminUpdateUserRole marks admin role changes.
	// AuditEventAdminUpdateUserRole 标记管理员修改用户角色事件。
	AuditEventAdminUpdateUserRole = "admin_update_user_role"
	// AuditEventAdminUpdateUserStatus marks admin status changes.
	// AuditEventAdminUpdateUserStatus 标记管理员修改用户状态事件。
	AuditEventAdminUpdateUserStatus = "admin_update_user_status"
	// AuditEventAdminResetUserPassword marks admin password resets.
	// AuditEventAdminResetUserPassword 标记管理员重置用户密码事件。
	AuditEventAdminResetUserPassword = "admin_reset_user_password"
	// AuditEventUpdateOwnProfile marks self-service profile updates.
	// AuditEventUpdateOwnProfile 标记用户更新自身资料事件。
	AuditEventUpdateOwnProfile = "update_own_profile"
	// AuditEventChangeOwnPassword marks self-service password changes.
	// AuditEventChangeOwnPassword 标记用户修改自身密码事件。
	AuditEventChangeOwnPassword = "change_own_password"
)
