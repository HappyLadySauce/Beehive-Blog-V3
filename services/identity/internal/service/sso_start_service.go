package service

import (
	"context"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
)

// SSOStartService handles outbound SSO authorize URL generation.
// SSOStartService 处理对外 SSO 授权地址生成。
type SSOStartService struct {
	deps Dependencies
}

// NewSSOStartService creates an SSOStartService instance.
// NewSSOStartService 创建 SSOStartService 实例。
func NewSSOStartService(deps Dependencies) *SSOStartService {
	return &SSOStartService{deps: deps}
}

// Execute validates provider readiness and persists OAuth state.
// Execute 校验 provider 就绪状态并持久化 OAuth state。
func (s *SSOStartService) Execute(ctx context.Context, in StartSSOInput) (*StartSSOResult, error) {
	providerName, err := auth.NormalizeProvider(in.Provider)
	if err != nil {
		return nil, errs.Wrap(err, errs.CodeIdentityInvalidArgument, "unsupported provider")
	}

	redirectURI := strings.TrimSpace(in.RedirectURI)
	if redirectURI == "" {
		return nil, errs.New(errs.CodeIdentityInvalidArgument, "redirect_uri is required")
	}

	providerItem, ok := s.deps.Providers.Get(providerName)
	if !ok {
		return nil, errs.New(errs.CodeIdentityInvalidArgument, "unsupported provider")
	}
	if !providerItem.Enabled() {
		return nil, errs.New(errs.CodeIdentitySSOProviderDisabled, "sso provider is disabled")
	}
	if redirectURI != strings.TrimSpace(providerItem.RedirectURL()) {
		return nil, errs.New(errs.CodeIdentityInvalidArgument, "redirect_uri does not match configured provider redirect")
	}
	if !providerItem.LoginReady() {
		writeAudit(ctx, s.deps.Store, auditInput{
			Provider:  stringPtr(providerName),
			EventType: auth.AuditEventStartSSO,
			Result:    auth.AuditResultFailure,
			ClientIP:  stringPtr(in.ClientIP),
			Detail: auth.MarshalAuditDetail(map[string]any{
				"reason": "sso_provider_not_ready",
			}),
		})
		return nil, errs.New(errs.CodeIdentitySSOProviderNotReady, "sso provider is not ready")
	}

	state := auth.EnsureState(in.State)
	authURL, err := providerItem.BuildAuthorizeURL(state)
	if err != nil {
		return nil, err
	}

	now := s.deps.Clock()
	stateRow := &entity.OAuthLoginState{
		Provider:    providerName,
		State:       state,
		RedirectURI: redirectURI,
		ExpiresAt:   now.Add(time.Duration(s.deps.Config.Security.StateTTLSeconds) * time.Second),
	}
	if err := s.deps.Store.OAuthLoginStates.Create(ctx, stateRow); err != nil {
		return nil, err
	}

	writeAudit(ctx, s.deps.Store, auditInput{
		Provider:  stringPtr(providerName),
		EventType: auth.AuditEventStartSSO,
		Result:    auth.AuditResultSuccess,
		ClientIP:  stringPtr(in.ClientIP),
	})

	return &StartSSOResult{
		Provider: providerName,
		AuthURL:  authURL,
		State:    state,
	}, nil
}
