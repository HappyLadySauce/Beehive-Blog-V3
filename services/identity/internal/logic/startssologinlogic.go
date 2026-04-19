package logic

import (
	"context"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/ctxmeta"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/svc"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/pb"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// StartSsoLoginLogic handles outbound SSO authorize URL generation.
// StartSsoLoginLogic 负责生成对外 SSO 授权地址。
type StartSsoLoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// NewStartSsoLoginLogic creates a StartSsoLoginLogic instance.
// NewStartSsoLoginLogic 创建 StartSsoLoginLogic 实例。
func NewStartSsoLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *StartSsoLoginLogic {
	return &StartSsoLoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// StartSsoLogin validates provider readiness and returns an authorize URL.
// StartSsoLogin 校验 provider 就绪状态并返回授权地址。
func (l *StartSsoLoginLogic) StartSsoLogin(in *pb.StartSsoLoginRequest) (*pb.StartSsoLoginResponse, error) {
	provider, err := auth.NormalizeProvider(in.GetProvider())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "unsupported provider")
	}
	redirectURI := strings.TrimSpace(in.GetRedirectUri())
	if redirectURI == "" {
		return nil, status.Error(codes.InvalidArgument, "redirect_uri is required")
	}

	if !auth.IsProviderEnabled(l.svcCtx.Config.SSO, provider) {
		return nil, status.Error(codes.FailedPrecondition, "sso_provider_disabled")
	}

	providerConf, err := auth.GetProviderConfig(l.svcCtx.Config.SSO, provider)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "unsupported provider")
	}
	if redirectURI != strings.TrimSpace(providerConf.RedirectURL) {
		return nil, status.Error(codes.InvalidArgument, "redirect_uri does not match configured provider redirect")
	}

	if !auth.IsProviderLoginReady(provider) {
		clientIP := ctxmeta.GetClientIPFromIncomingContext(l.ctx)
		writeAudit(l.ctx, l.svcCtx.Store, auditInput{
			Provider:  stringPtr(provider),
			EventType: auth.AuditEventStartSSO,
			Result:    auth.AuditResultFailure,
			ClientIP:  stringPtr(clientIP),
			Detail: auth.MarshalAuditDetail(map[string]any{
				"reason": "sso_provider_not_ready",
			}),
		})
		l.Infof("sso start rejected because provider is not ready: provider=%s", provider)
		return nil, status.Error(codes.FailedPrecondition, "sso_provider_not_ready")
	}

	state := auth.EnsureState(in.GetState())
	authURL, err := auth.BuildAuthorizeURL(provider, providerConf, state)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "build authorize URL failed: %v", err)
	}

	now := time.Now().UTC()
	stateRow := &entity.OAuthLoginState{
		Provider:        provider,
		State:           state,
		RedirectURI:     redirectURI,
		RequestedScopes: stringPtr(strings.Join(providerConf.Scopes, ",")),
		ExpiresAt:       now.Add(time.Duration(l.svcCtx.Config.Security.StateTTLSeconds) * time.Second),
	}
	if err := l.svcCtx.Store.OAuthLoginStates.Create(l.ctx, stateRow); err != nil {
		return nil, status.Errorf(codes.Internal, "persist oauth state failed: %v", err)
	}

	writeAudit(l.ctx, l.svcCtx.Store, auditInput{
		Provider:  stringPtr(provider),
		EventType: auth.AuditEventStartSSO,
		Result:    auth.AuditResultSuccess,
		ClientIP:  stringPtr(ctxmeta.GetClientIPFromIncomingContext(l.ctx)),
	})
	l.Infof("sso start succeeded: provider=%s", provider)

	return &pb.StartSsoLoginResponse{
		Provider: provider,
		AuthUrl:  authURL,
		State:    state,
	}, nil
}
