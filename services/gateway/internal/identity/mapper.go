package identity

import (
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
)

// BuildRegisterRequest maps HTTP register request to identity RPC request.
// BuildRegisterRequest 将 HTTP 注册请求转换为 identity RPC 请求。
func BuildRegisterRequest(req *types.AuthRegisterReq) *pb.RegisterLocalUserRequest {
	return &pb.RegisterLocalUserRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Nickname: req.Nickname,
	}
}

// BuildLoginRequest maps HTTP login request to identity RPC request.
// BuildLoginRequest 将 HTTP 登录请求转换为 identity RPC 请求。
func BuildLoginRequest(req *types.AuthLoginReq) *pb.LoginLocalUserRequest {
	return &pb.LoginLocalUserRequest{
		LoginIdentifier: req.LoginIdentifier,
		Password:        req.Password,
		ClientType:      req.ClientType,
		DeviceId:        req.DeviceId,
		DeviceName:      req.DeviceName,
		UserAgent:       req.UserAgent,
	}
}

// BuildSsoStartRequest maps HTTP SSO start request.
// BuildSsoStartRequest 转换 SSO 发起请求。
func BuildSsoStartRequest(req *types.AuthSsoStartReq) *pb.StartSsoLoginRequest {
	return &pb.StartSsoLoginRequest{
		Provider:    req.Provider,
		RedirectUri: req.RedirectUri,
		State:       req.State,
	}
}

// BuildSsoCallbackRequest maps HTTP SSO callback request.
// BuildSsoCallbackRequest 转换 SSO 回调请求。
func BuildSsoCallbackRequest(req *types.AuthSsoCallbackReq) *pb.FinishSsoLoginRequest {
	return &pb.FinishSsoLoginRequest{
		Provider:    req.Provider,
		Code:        req.Code,
		State:       req.State,
		RedirectUri: req.RedirectUri,
		ClientType:  req.ClientType,
		DeviceId:    req.DeviceId,
		DeviceName:  req.DeviceName,
		UserAgent:   req.UserAgent,
	}
}

// BuildRefreshRequest maps HTTP refresh request.
// BuildRefreshRequest 转换刷新 token 请求。
func BuildRefreshRequest(req *types.AuthRefreshReq) *pb.RefreshSessionTokenRequest {
	return &pb.RefreshSessionTokenRequest{
		RefreshToken: req.RefreshToken,
		UserAgent:    req.UserAgent,
	}
}

// BuildLogoutRequest maps trusted logout context and body to RPC request.
// BuildLogoutRequest 将可信上下文与请求体转换为登出 RPC 请求。
func BuildLogoutRequest(sessionID string, req *types.AuthLogoutReq) *pb.LogoutSessionRequest {
	return &pb.LogoutSessionRequest{
		SessionId:    sessionID,
		RefreshToken: req.RefreshToken,
	}
}

// BuildMeRequest maps trusted user id to RPC request.
// BuildMeRequest 将可信 user id 转换为 RPC 请求。
func BuildMeRequest(userID string) *pb.GetCurrentUserRequest {
	return &pb.GetCurrentUserRequest{
		UserId: userID,
	}
}

// BuildUpdateProfileRequest maps the current-user profile patch request.
// BuildUpdateProfileRequest 转换当前用户资料更新请求。
func BuildUpdateProfileRequest(userID string, req *types.AuthUpdateProfileReq) *pb.UpdateOwnProfileRequest {
	return &pb.UpdateOwnProfileRequest{
		UserId:    userID,
		Nickname:  req.Nickname,
		AvatarUrl: req.AvatarUrl,
	}
}

// BuildChangePasswordRequest maps the current-user password change request.
// BuildChangePasswordRequest 转换当前用户密码修改请求。
func BuildChangePasswordRequest(userID string, req *types.AuthChangePasswordReq) *pb.ChangeOwnPasswordRequest {
	return &pb.ChangeOwnPasswordRequest{
		UserId:      userID,
		OldPassword: req.OldPassword,
		NewPassword: req.NewPassword,
	}
}

// BuildListUsersRequest maps studio user list filters.
// BuildListUsersRequest 转换 Studio 用户列表过滤条件。
func BuildListUsersRequest(actorUserID string, req *types.AdminUserListReq) *pb.ListUsersRequest {
	return &pb.ListUsersRequest{
		ActorUserId: actorUserID,
		Keyword:     req.Keyword,
		Role:        toProtoRole(req.Role),
		Status:      toProtoAccountStatus(req.Status),
		Page:        int32(req.Page),
		PageSize:    int32(req.PageSize),
	}
}

// BuildUpdateUserRoleRequest maps studio user role updates.
// BuildUpdateUserRoleRequest 转换 Studio 用户角色更新请求。
func BuildUpdateUserRoleRequest(actorUserID string, req *types.AdminUpdateUserRoleReq) *pb.UpdateUserRoleRequest {
	return &pb.UpdateUserRoleRequest{
		ActorUserId:  actorUserID,
		TargetUserId: req.UserId,
		Role:         toProtoRole(req.Role),
	}
}

// BuildUpdateUserStatusRequest maps studio user status updates.
// BuildUpdateUserStatusRequest 转换 Studio 用户状态更新请求。
func BuildUpdateUserStatusRequest(actorUserID string, req *types.AdminUpdateUserStatusReq) *pb.UpdateUserStatusRequest {
	return &pb.UpdateUserStatusRequest{
		ActorUserId:  actorUserID,
		TargetUserId: req.UserId,
		Status:       toProtoAccountStatus(req.Status),
	}
}

// BuildResetUserPasswordRequest maps studio password reset requests.
// BuildResetUserPasswordRequest 转换 Studio 用户密码重置请求。
func BuildResetUserPasswordRequest(actorUserID string, req *types.AdminResetUserPasswordReq) *pb.ResetUserPasswordRequest {
	return &pb.ResetUserPasswordRequest{
		ActorUserId:  actorUserID,
		TargetUserId: req.UserId,
		NewPassword:  req.NewPassword,
	}
}

// BuildListAuditsRequest maps identity audit filters.
// BuildListAuditsRequest 转换 identity 审计列表过滤条件。
func BuildListAuditsRequest(actorUserID string, req *types.IdentityAuditListReq) *pb.ListIdentityAuditsRequest {
	return &pb.ListIdentityAuditsRequest{
		ActorUserId: actorUserID,
		EventType:   req.EventType,
		Result:      req.Result,
		UserId:      req.UserId,
		StartedAt:   req.StartedAt,
		EndedAt:     req.EndedAt,
		Page:        int32(req.Page),
		PageSize:    int32(req.PageSize),
	}
}

// ToRegisterResponse maps identity register response to HTTP response.
// ToRegisterResponse 将 identity 注册响应转换为 HTTP 响应。
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
		AuthSource: session.GetAuthSource().String(),
		ClientType: session.GetClientType(),
		DeviceId:   session.GetDeviceId(),
		DeviceName: session.GetDeviceName(),
		Status:     session.GetStatus().String(),
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

func toProtoRole(role string) pb.Role {
	switch strings.ToLower(strings.TrimSpace(role)) {
	case "admin", "role_admin":
		return pb.Role_ROLE_ADMIN
	case "member", "user", "role_member":
		return pb.Role_ROLE_MEMBER
	case "guest", "role_guest":
		return pb.Role_ROLE_GUEST
	default:
		return pb.Role_ROLE_UNSPECIFIED
	}
}

func fromProtoRole(role pb.Role) string {
	switch role {
	case pb.Role_ROLE_ADMIN:
		return "admin"
	case pb.Role_ROLE_MEMBER:
		return "member"
	case pb.Role_ROLE_GUEST:
		return "guest"
	default:
		return ""
	}
}

func toProtoAccountStatus(status string) pb.AccountStatus {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "pending", "account_status_pending":
		return pb.AccountStatus_ACCOUNT_STATUS_PENDING
	case "active", "account_status_active":
		return pb.AccountStatus_ACCOUNT_STATUS_ACTIVE
	case "disabled", "account_status_disabled":
		return pb.AccountStatus_ACCOUNT_STATUS_DISABLED
	case "locked", "account_status_locked":
		return pb.AccountStatus_ACCOUNT_STATUS_LOCKED
	default:
		return pb.AccountStatus_ACCOUNT_STATUS_UNSPECIFIED
	}
}

func fromProtoAccountStatus(status pb.AccountStatus) string {
	switch status {
	case pb.AccountStatus_ACCOUNT_STATUS_PENDING:
		return "pending"
	case pb.AccountStatus_ACCOUNT_STATUS_ACTIVE:
		return "active"
	case pb.AccountStatus_ACCOUNT_STATUS_DISABLED:
		return "disabled"
	case pb.AccountStatus_ACCOUNT_STATUS_LOCKED:
		return "locked"
	default:
		return ""
	}
}

func fromProtoAuthSource(authSource pb.AuthSource) string {
	switch authSource {
	case pb.AuthSource_AUTH_SOURCE_LOCAL:
		return "local"
	case pb.AuthSource_AUTH_SOURCE_SSO:
		return "sso"
	default:
		return ""
	}
}
