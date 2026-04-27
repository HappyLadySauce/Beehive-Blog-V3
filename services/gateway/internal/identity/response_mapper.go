package identity

import (
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
)

func ToRegisterResponse(resp *pb.RegisterLocalUserResponse) *types.AuthRegisterResp {
	return &types.AuthRegisterResp{
		AccessToken:  safeTokenPair(resp.GetTokenPair()).AccessToken,
		RefreshToken: safeTokenPair(resp.GetTokenPair()).RefreshToken,
		ExpiresIn:    safeTokenPair(resp.GetTokenPair()).ExpiresIn,
		TokenType:    safeTokenPair(resp.GetTokenPair()).TokenType,
		SessionId:    safeSession(resp.GetSessionInfo()).SessionId,
		User:         toUserProfile(resp.GetCurrentUser()),
		Session:      toSessionView(resp.GetSessionInfo()),
	}
}

// ToLoginResponse maps identity login response to HTTP response.
// ToLoginResponse 将 identity 登录响应转换为 HTTP 响应。

func ToLoginResponse(resp *pb.LoginLocalUserResponse) *types.AuthLoginResp {
	return &types.AuthLoginResp{
		AccessToken:  safeTokenPair(resp.GetTokenPair()).AccessToken,
		RefreshToken: safeTokenPair(resp.GetTokenPair()).RefreshToken,
		ExpiresIn:    safeTokenPair(resp.GetTokenPair()).ExpiresIn,
		TokenType:    safeTokenPair(resp.GetTokenPair()).TokenType,
		SessionId:    safeSession(resp.GetSessionInfo()).SessionId,
		User:         toUserProfile(resp.GetCurrentUser()),
		Session:      toSessionView(resp.GetSessionInfo()),
	}
}

// ToSsoStartResponse maps identity SSO start response.
// ToSsoStartResponse 转换 identity SSO 发起响应。

func ToSsoStartResponse(resp *pb.StartSsoLoginResponse) *types.AuthSsoStartResp {
	return &types.AuthSsoStartResp{
		Provider: resp.GetProvider(),
		AuthUrl:  resp.GetAuthUrl(),
		State:    resp.GetState(),
	}
}

// ToSsoCallbackResponse maps identity SSO callback response.
// ToSsoCallbackResponse 转换 identity SSO 回调响应。

func ToSsoCallbackResponse(resp *pb.FinishSsoLoginResponse) *types.AuthSsoCallbackResp {
	return &types.AuthSsoCallbackResp{
		AccessToken:  safeTokenPair(resp.GetTokenPair()).AccessToken,
		RefreshToken: safeTokenPair(resp.GetTokenPair()).RefreshToken,
		ExpiresIn:    safeTokenPair(resp.GetTokenPair()).ExpiresIn,
		TokenType:    safeTokenPair(resp.GetTokenPair()).TokenType,
		SessionId:    safeSession(resp.GetSessionInfo()).SessionId,
		User:         toUserProfile(resp.GetCurrentUser()),
		Session:      toSessionView(resp.GetSessionInfo()),
	}
}

// ToRefreshResponse maps identity refresh response.
// ToRefreshResponse 转换 identity 刷新响应。

func ToRefreshResponse(resp *pb.RefreshSessionTokenResponse) *types.AuthRefreshResp {
	return &types.AuthRefreshResp{
		AccessToken:  safeTokenPair(resp.GetTokenPair()).AccessToken,
		RefreshToken: safeTokenPair(resp.GetTokenPair()).RefreshToken,
		ExpiresIn:    safeTokenPair(resp.GetTokenPair()).ExpiresIn,
		TokenType:    safeTokenPair(resp.GetTokenPair()).TokenType,
		SessionId:    safeSession(resp.GetSessionInfo()).SessionId,
		Session:      toSessionView(resp.GetSessionInfo()),
	}
}

// ToMeResponse maps identity current user response.
// ToMeResponse 转换 identity 当前用户响应。

func ToMeResponse(resp *pb.GetCurrentUserResponse) *types.AuthMeResp {
	return &types.AuthMeResp{
		User: toUserProfile(resp.GetCurrentUser()),
	}
}

// ToUpdateProfileResponse maps identity profile update response.
// ToUpdateProfileResponse 转换 identity 资料更新响应。

func ToUpdateProfileResponse(resp *pb.UpdateOwnProfileResponse) *types.AuthMeResp {
	return &types.AuthMeResp{
		User: toUserProfile(resp.GetCurrentUser()),
	}
}

// ToChangePasswordResponse maps identity password change response.
// ToChangePasswordResponse 转换 identity 密码修改响应。

func ToChangePasswordResponse(resp *pb.ChangeOwnPasswordResponse) *types.AuthChangePasswordResp {
	return &types.AuthChangePasswordResp{
		Ok: resp.GetOk(),
	}
}

// ToAdminUserListResponse maps identity users list response.
// ToAdminUserListResponse 转换 identity 用户列表响应。

func ToAdminUserListResponse(resp *pb.ListUsersResponse) *types.AdminUserListResp {
	items := make([]types.AdminUserView, 0, len(resp.GetItems()))
	for _, item := range resp.GetItems() {
		items = append(items, toAdminUserView(item))
	}
	return &types.AdminUserListResp{
		Items:    items,
		Total:    resp.GetTotal(),
		Page:     int(resp.GetPage()),
		PageSize: int(resp.GetPageSize()),
	}
}

// ToAdminUserResponse maps identity admin user response.
// ToAdminUserResponse 转换 identity 管理用户响应。

func ToAdminUserResponse(user *pb.AdminUserView) *types.AdminUserResp {
	return &types.AdminUserResp{
		User: toAdminUserView(user),
	}
}

// ToResetUserPasswordResponse maps identity password reset response.
// ToResetUserPasswordResponse 转换 identity 密码重置响应。

func ToResetUserPasswordResponse(resp *pb.ResetUserPasswordResponse) *types.AdminResetUserPasswordResp {
	return &types.AdminResetUserPasswordResp{
		Ok: resp.GetOk(),
	}
}

// ToAuditListResponse maps identity audit list response.
// ToAuditListResponse 转换 identity 审计列表响应。

func ToAuditListResponse(resp *pb.ListIdentityAuditsResponse) *types.IdentityAuditListResp {
	items := make([]types.IdentityAuditView, 0, len(resp.GetItems()))
	for _, item := range resp.GetItems() {
		items = append(items, toAuditView(item))
	}
	return &types.IdentityAuditListResp{
		Items:    items,
		Total:    resp.GetTotal(),
		Page:     int(resp.GetPage()),
		PageSize: int(resp.GetPageSize()),
	}
}

func toUserProfile(user *pb.CurrentUser) types.AuthUserProfile {
	if user == nil {
		return types.AuthUserProfile{}
	}
	return types.AuthUserProfile{
		UserId:    user.GetUserId(),
		Username:  user.GetUsername(),
		Email:     user.GetEmail(),
		Nickname:  user.GetNickname(),
		AvatarUrl: user.GetAvatarUrl(),
		Role:      fromProtoRole(user.GetRole()),
		Status:    fromProtoAccountStatus(user.GetStatus()),
	}
}

func toAdminUserView(user *pb.AdminUserView) types.AdminUserView {
	if user == nil {
		return types.AdminUserView{}
	}
	return types.AdminUserView{
		UserId:      user.GetUserId(),
		Username:    user.GetUsername(),
		Email:       user.GetEmail(),
		Nickname:    user.GetNickname(),
		AvatarUrl:   user.GetAvatarUrl(),
		Role:        fromProtoRole(user.GetRole()),
		Status:      fromProtoAccountStatus(user.GetStatus()),
		LastLoginAt: user.GetLastLoginAt(),
		CreatedAt:   user.GetCreatedAt(),
		UpdatedAt:   user.GetUpdatedAt(),
	}
}

func toAuditView(audit *pb.IdentityAuditView) types.IdentityAuditView {
	if audit == nil {
		return types.IdentityAuditView{}
	}
	return types.IdentityAuditView{
		AuditId:    audit.GetAuditId(),
		UserId:     audit.GetUserId(),
		SessionId:  audit.GetSessionId(),
		Provider:   audit.GetProvider(),
		AuthSource: fromProtoAuthSource(audit.GetAuthSource()),
		EventType:  audit.GetEventType(),
		Result:     audit.GetResult(),
		ClientIp:   audit.GetClientIp(),
		UserAgent:  audit.GetUserAgent(),
		DetailJson: audit.GetDetailJson(),
		CreatedAt:  audit.GetCreatedAt(),
	}
}

func toSessionView(session *pb.SessionInfo) types.AuthSessionView {
	if session == nil {
		return types.AuthSessionView{}
	}
	return types.AuthSessionView{
		SessionId:  session.GetSessionId(),
		UserId:     session.GetUserId(),
		AuthSource: fromProtoAuthSource(session.GetAuthSource()),
		ClientType: session.GetClientType(),
		DeviceId:   session.GetDeviceId(),
		DeviceName: session.GetDeviceName(),
		Status:     fromProtoSessionStatus(session.GetStatus()),
		LastSeenAt: session.GetLastSeenAt(),
		ExpiresAt:  session.GetExpiresAt(),
	}
}

func safeTokenPair(tokenPair *pb.TokenPair) *pb.TokenPair {
	if tokenPair != nil {
		return tokenPair
	}
	return &pb.TokenPair{}
}

func safeSession(session *pb.SessionInfo) *pb.SessionInfo {
	if session != nil {
		return session
	}
	return &pb.SessionInfo{}
}
