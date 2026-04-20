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
	ExchangeCode(ctx context.Context, code, redirectURI string) (string, error)
	FetchProfile(ctx context.Context, accessToken string) (*Profile, []byte, error)
}

// Profile is the normalized federated profile produced by callback providers.
// Profile 是回调 provider 生成的标准化联邦资料。
type Profile struct {
	Subject          string
	SubjectType      string
	Login            string
	DisplayName      string
	Email            *string
	AvatarURL        *string
	RawProfile       []byte
	ProviderClientID *string
	RequestedScopes  *string
}
