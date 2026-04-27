package identity

import (
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/gateway/internal/types"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
)

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
