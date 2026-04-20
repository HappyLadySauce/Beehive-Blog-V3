package identity

import (
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
		Role:      user.GetRole().String(),
		Status:    user.GetStatus().String(),
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
