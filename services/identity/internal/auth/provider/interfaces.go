package provider

import "context"

// Provider defines the minimal capabilities required for outbound SSO start flow.
// Provider 定义对外 SSO start 流程所需的最小能力集合。
type Provider interface {
	Name() string
	Enabled() bool
	LoginReady() bool
	RedirectURL() string
	BuildAuthorizeURL(state string) (string, error)
}

// CallbackProvider extends Provider with callback completion capabilities.
// CallbackProvider 在 Provider 基础上扩展回调完成能力。
type CallbackProvider interface {
	Provider
	ExchangeCode(ctx context.Context, code, redirectURI string) (*AccessToken, error)
	FetchProfile(ctx context.Context, accessToken *AccessToken) (*Profile, []byte, error)
}

// AccessToken contains the normalized provider token exchange result.
// AccessToken 包含规范化后的 provider token 交换结果。
type AccessToken struct {
	Token        string
	RefreshToken string
	OpenID       string
	UnionID      string
	Scope        string
}

// Profile is the normalized federated profile produced by callback providers.
// Profile 是回调 provider 生成的标准化联邦资料。
type Profile struct {
	Subject          string
	SubjectType      string
	Login            string
	DisplayName      string
	Email            *string
	EmailVerified    bool
	AvatarURL        *string
	UnionID          *string
	OpenID           *string
	RawProfile       []byte
	ProviderClientID *string
	RequestedScopes  *string
}

var _ CallbackProvider = (*GitHubClient)(nil)
var _ CallbackProvider = (*QQClient)(nil)
var _ CallbackProvider = (*WeChatClient)(nil)
