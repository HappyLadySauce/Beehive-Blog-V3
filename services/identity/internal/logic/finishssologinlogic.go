package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	identityservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
	"github.com/zeromicro/go-zero/core/logx"
)

// FinishSsoLoginLogic handles SSO callback completion.
// FinishSsoLoginLogic 负责处理 SSO callback 完成。
type FinishSsoLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewFinishSsoLoginLogic creates a FinishSsoLoginLogic instance.
// NewFinishSsoLoginLogic 创建 FinishSsoLoginLogic 实例。
func NewFinishSsoLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *FinishSsoLoginLogic {
	return &FinishSsoLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// FinishSsoLogin adapts SSO callback completion to the service layer.
// FinishSsoLogin 将 SSO callback 完成流程适配到 service 层。
func (l *FinishSsoLoginLogic) FinishSsoLogin(in *pb.FinishSsoLoginRequest) (*pb.FinishSsoLoginResponse, error) {
	result, err := l.svcCtx.Services.SSOFinish.Execute(l.ctx, identityservice.FinishSSOInput{
		Provider:    in.GetProvider(),
		Code:        in.GetCode(),
		State:       in.GetState(),
		RedirectURI: in.GetRedirectUri(),
		ClientType:  in.GetClientType(),
		DeviceID:    in.GetDeviceId(),
		DeviceName:  in.GetDeviceName(),
		UserAgent:   in.GetUserAgent(),
		ClientIP:    ctxmeta.GetClientIPFromIncomingContext(l.ctx),
	})
	if err != nil {
		return nil, toStatusError(err, "finish sso login failed")
	}

	l.Infof("sso finish succeeded: user_id=%d session_id=%d", result.User.ID, result.Session.ID)

	return &pb.FinishSsoLoginResponse{
		TokenPair: auth.NewTokenPair(
			result.AccessToken,
			result.RefreshToken,
			expiresIn(result.AccessExpiresAt),
			result.Session.ID,
		),
		CurrentUser: auth.ToCurrentUser(result.User),
		SessionInfo: auth.ToSessionInfo(result.Session),
	}, nil
}
