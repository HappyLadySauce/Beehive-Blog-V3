package logic

import (
	"context"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/logs"
	identityservice "github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/service"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
)

// StartSsoLoginLogic handles outbound SSO authorize URL generation.
// StartSsoLoginLogic 负责生成对外 SSO 授权地址。
type StartSsoLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logger *logs.Logger
}

// NewStartSsoLoginLogic creates a StartSsoLoginLogic instance.
// NewStartSsoLoginLogic 创建 StartSsoLoginLogic 实例。
func NewStartSsoLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartSsoLoginLogic {
	logCtx := withLogContext(ctx)
	return &StartSsoLoginLogic{
		ctx:    logCtx,
		svcCtx: svcCtx,
		logger: logs.Ctx(logCtx),
	}
}

// StartSsoLogin adapts SSO start requests to the service layer.
// StartSsoLogin 将 SSO start 请求适配到 service 层。
func (l *StartSsoLoginLogic) StartSsoLogin(in *pb.StartSsoLoginRequest) (*pb.StartSsoLoginResponse, error) {
	result, err := l.svcCtx.Services.SSOStart.Execute(l.ctx, identityservice.StartSSOInput{
		Provider:    in.GetProvider(),
		RedirectURI: in.GetRedirectUri(),
		State:       in.GetState(),
		ClientIP:    ctxmeta.GetClientIPFromIncomingContext(l.ctx),
	})
	if err != nil {
		return nil, toStatusError(err, "start sso login failed")
	}

	l.logger.Info(
		"identity_start_sso_login_succeeded",
		logs.String("provider", result.Provider),
	)

	return &pb.StartSsoLoginResponse{
		Provider: result.Provider,
		AuthUrl:  result.AuthURL,
		State:    result.State,
	}, nil
}
