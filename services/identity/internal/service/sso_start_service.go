package service

import (
	"context"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/pkg/errs"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/entity"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/model/repo"
)

const (
	oauthStatePurposeLogin       = "login"
	oauthStatePurposeEmailUpdate = "email_update"
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
	return s.start(ctx, ssoStartParams{
		Provider:    in.Provider,
		RedirectURI: in.RedirectURI,
		State:       in.State,
		ClientIP:    in.ClientIP,
		Purpose:     oauthStatePurposeLogin,
		EventType:   auth.AuditEventStartSSO,
	})
}

// ExecuteReauth starts an SSO reauthentication flow for sensitive account changes.
// ExecuteReauth 为敏感账号变更发起 SSO 重验流程。
func (s *SSOStartService) ExecuteReauth(ctx context.Context, in StartSSOReauthInput) (*StartSSOResult, error) {
	if in.UserID <= 0 {
		return nil, errs.New(errs.CodeIdentityInvalidArgument, "user_id is invalid")
	}
	user, err := s.deps.Store.Users.GetByID(ctx, in.UserID)
	if err != nil {
		if repo.IsNotFound(err) {
			return nil, errs.New(errs.CodeIdentityUserNotFound, "user not found")
		}
		return nil, errs.Wrap(err, errs.CodeIdentityInternal, "load user failed")
	}
	if err := validateActiveUserStatus(user.Status); err != nil {
		return nil, err
	}

	return s.start(ctx, ssoStartParams{
		Provider:      in.Provider,
		RedirectURI:   in.RedirectURI,
		State:         in.State,
		ClientIP:      in.ClientIP,
		Purpose:       oauthStatePurposeEmailUpdate,
		SubjectUserID: &in.UserID,
		EventType:     auth.AuditEventStartSSOReauth,
	})
}

type ssoStartParams struct {
	Provider      string
	RedirectURI   string
	State         string
	ClientIP      string
	Purpose       string
	SubjectUserID *int64
	EventType     string
}

func (s *SSOStartService) start(ctx context.Context, in ssoStartParams) (*StartSSOResult, error) {
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
			EventType: in.EventType,
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
		Provider:      providerName,
		State:         state,
		RedirectURI:   redirectURI,
		Purpose:       in.Purpose,
		SubjectUserID: in.SubjectUserID,
		ExpiresAt:     now.Add(time.Duration(s.deps.Config.Security.StateTTLSeconds) * time.Second),
	}
	if err := s.deps.Store.OAuthLoginStates.Create(ctx, stateRow); err != nil {
		return nil, err
	}

	writeAudit(ctx, s.deps.Store, auditInput{
		Provider:  stringPtr(providerName),
		EventType: in.EventType,
		Result:    auth.AuditResultSuccess,
		ClientIP:  stringPtr(in.ClientIP),
	})

	return &StartSSOResult{
		Provider: providerName,
		AuthURL:  authURL,
		State:    state,
	}, nil
}
