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
