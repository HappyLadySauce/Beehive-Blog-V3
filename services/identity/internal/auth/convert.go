package auth

import (
	"strconv"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
)

// ToCurrentUser converts a user entity to proto.
// ToCurrentUser 将用户实体转换为 proto。
func ToCurrentUser(user *entity.User) *pb.CurrentUser {
	if user == nil {
		return nil
	}

	return &pb.CurrentUser{
		UserId:    strconv.FormatInt(user.ID, 10),
		Username:  user.Username,
		Email:     derefString(user.Email),
		Nickname:  derefString(user.Nickname),
		AvatarUrl: derefString(user.AvatarURL),
		Role:      ToProtoRole(user.Role),
		Status:    ToProtoAccountStatus(user.Status),
	}
}

// ToAdminUserView converts a user entity to admin proto view.
// ToAdminUserView 将用户实体转换为管理员视图 proto。
func ToAdminUserView(user *entity.User) *pb.AdminUserView {
	if user == nil {
		return nil
	}

	var lastLoginAt int64
	if user.LastLoginAt != nil {
		lastLoginAt = user.LastLoginAt.Unix()
	}

	return &pb.AdminUserView{
		UserId:      strconv.FormatInt(user.ID, 10),
		Username:    user.Username,
		Email:       derefString(user.Email),
		Nickname:    derefString(user.Nickname),
		AvatarUrl:   derefString(user.AvatarURL),
		Role:        ToProtoRole(user.Role),
		Status:      ToProtoAccountStatus(user.Status),
		LastLoginAt: lastLoginAt,
		CreatedAt:   user.CreatedAt.Unix(),
		UpdatedAt:   user.UpdatedAt.Unix(),
	}
}

// ToIdentityAuditView converts an audit entity to proto view.
// ToIdentityAuditView 将审计实体转换为 proto 视图。
func ToIdentityAuditView(audit *entity.IdentityAudit) *pb.IdentityAuditView {
	if audit == nil {
		return nil
	}

	return &pb.IdentityAuditView{
		AuditId:    strconv.FormatInt(audit.ID, 10),
		UserId:     int64PtrString(audit.UserID),
		SessionId:  int64PtrString(audit.SessionID),
		Provider:   derefString(audit.Provider),
		AuthSource: ToProtoAuthSource(derefString(audit.AuthSource)),
		EventType:  audit.EventType,
		Result:     audit.Result,
		ClientIp:   derefString(audit.ClientIP),
		UserAgent:  derefString(audit.UserAgent),
		DetailJson: string(audit.Detail),
		CreatedAt:  audit.CreatedAt.Unix(),
	}
}

// ToSessionInfo converts a session entity to proto.
// ToSessionInfo 将会话实体转换为 proto。
func ToSessionInfo(session *entity.UserSession) *pb.SessionInfo {
	if session == nil {
		return nil
	}

	var lastSeenAt int64
	if session.LastSeenAt != nil {
		lastSeenAt = session.LastSeenAt.Unix()
	}

	return &pb.SessionInfo{
		SessionId:  strconv.FormatInt(session.ID, 10),
		UserId:     strconv.FormatInt(session.UserID, 10),
		AuthSource: ToProtoAuthSource(session.AuthSource),
		ClientType: derefString(session.ClientType),
		DeviceId:   derefString(session.DeviceID),
		DeviceName: derefString(session.DeviceName),
		Status:     ToProtoSessionStatus(session.Status),
		LastSeenAt: lastSeenAt,
		ExpiresAt:  session.ExpiresAt.Unix(),
	}
}

// NewTokenPair builds a token pair response.
// NewTokenPair 构建 token pair 响应。
func NewTokenPair(accessToken, refreshToken string, expiresIn int64, sessionID int64) *pb.TokenPair {
	return &pb.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    expiresIn,
		TokenType:    "Bearer",
		SessionId:    strconv.FormatInt(sessionID, 10),
	}
}

// ToProtoRole maps a role string to proto.
// ToProtoRole 将角色字符串映射到 proto。
func ToProtoRole(role string) pb.Role {
	switch role {
	case UserRoleAdmin:
		return pb.Role_ROLE_ADMIN
	case UserRoleMember:
		return pb.Role_ROLE_MEMBER
	default:
		return pb.Role_ROLE_UNSPECIFIED
	}
}

// FromProtoRole maps a role proto enum to a storage value.
// FromProtoRole 将角色 proto 枚举映射为存储值。
func FromProtoRole(role pb.Role) string {
	switch role {
	case pb.Role_ROLE_ADMIN:
		return UserRoleAdmin
	case pb.Role_ROLE_MEMBER:
		return UserRoleMember
	default:
		return ""
	}
}

// FromProtoAccountStatus maps an account-status proto enum to a storage value.
// FromProtoAccountStatus 将账号状态 proto 枚举映射为存储值。
func FromProtoAccountStatus(status pb.AccountStatus) string {
	switch status {
	case pb.AccountStatus_ACCOUNT_STATUS_PENDING:
		return UserStatusPending
	case pb.AccountStatus_ACCOUNT_STATUS_ACTIVE:
		return UserStatusActive
	case pb.AccountStatus_ACCOUNT_STATUS_DISABLED:
		return UserStatusDisabled
	case pb.AccountStatus_ACCOUNT_STATUS_LOCKED:
		return UserStatusLocked
	default:
		return ""
	}
}

// ToProtoAccountStatus maps an account status string to proto.
// ToProtoAccountStatus 将账号状态字符串映射到 proto。
func ToProtoAccountStatus(status string) pb.AccountStatus {
	switch status {
	case UserStatusPending:
		return pb.AccountStatus_ACCOUNT_STATUS_PENDING
	case UserStatusActive:
		return pb.AccountStatus_ACCOUNT_STATUS_ACTIVE
	case UserStatusDisabled:
		return pb.AccountStatus_ACCOUNT_STATUS_DISABLED
	case UserStatusLocked:
		return pb.AccountStatus_ACCOUNT_STATUS_LOCKED
	default:
		return pb.AccountStatus_ACCOUNT_STATUS_UNSPECIFIED
	}
}

// ToProtoAuthSource maps an auth source string to proto.
// ToProtoAuthSource 将认证来源字符串映射到 proto。
func ToProtoAuthSource(source string) pb.AuthSource {
	switch source {
	case AuthSourceLocal:
		return pb.AuthSource_AUTH_SOURCE_LOCAL
	case AuthSourceSSO:
		return pb.AuthSource_AUTH_SOURCE_SSO
	default:
		return pb.AuthSource_AUTH_SOURCE_UNSPECIFIED
	}
}

// ToProtoSessionStatus maps a session status string to proto.
// ToProtoSessionStatus 将会话状态字符串映射到 proto。
func ToProtoSessionStatus(status string) pb.SessionStatus {
	switch status {
	case SessionStatusActive:
		return pb.SessionStatus_SESSION_STATUS_ACTIVE
	case SessionStatusRevoked:
		return pb.SessionStatus_SESSION_STATUS_REVOKED
	case SessionStatusExpired:
		return pb.SessionStatus_SESSION_STATUS_EXPIRED
	default:
		return pb.SessionStatus_SESSION_STATUS_UNSPECIFIED
	}
}

func derefString(value *string) string {
	if value == nil {
		return ""
	}

	return *value
}

func int64PtrString(value *int64) string {
	if value == nil {
		return ""
	}
	return strconv.FormatInt(*value, 10)
}
