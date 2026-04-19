package auth

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/config"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	githuboauth "golang.org/x/oauth2/github"
)

// IsProviderLoginReady reports whether the provider has a complete login implementation.
// IsProviderLoginReady 判断 provider 是否具备完整登录实现。
func IsProviderLoginReady(provider string) bool {
	return strings.EqualFold(provider, ProviderGitHub)
}

// IsProviderEnabled reports whether the provider is configured and enabled.
// IsProviderEnabled 判断 provider 是否已配置并启用。
func IsProviderEnabled(conf config.SSOConf, provider string) bool {
	switch strings.ToLower(provider) {
	case ProviderGitHub:
		return conf.GitHub.Enabled
	case ProviderQQ:
		return conf.QQ.Enabled
	case ProviderWeChat:
		return conf.WeChat.Enabled
	default:
		return false
	}
}

// GetProviderConfig returns the provider config by provider key.
// GetProviderConfig 按 provider key 返回 provider 配置。
func GetProviderConfig(conf config.SSOConf, provider string) (config.OAuthProviderConf, error) {
	switch strings.ToLower(provider) {
	case ProviderGitHub:
		return conf.GitHub, nil
	case ProviderQQ:
		return conf.QQ, nil
	case ProviderWeChat:
		return conf.WeChat, nil
	default:
		return config.OAuthProviderConf{}, fmt.Errorf("unsupported provider")
	}
}

// EnsureState returns the provided state or generates a new one when it is empty.
// EnsureState 返回传入 state，若为空则自动生成。
func EnsureState(state string) string {
	if trimmed := strings.TrimSpace(state); trimmed != "" {
		return trimmed
	}

	return uuid.NewString()
}

// BuildAuthorizeURL builds a provider-specific authorize URL.
// BuildAuthorizeURL 构建 provider 专属授权地址。
func BuildAuthorizeURL(provider string, conf config.OAuthProviderConf, state string) (string, error) {
	switch strings.ToLower(provider) {
	case ProviderGitHub:
		oauthConf := oauth2.Config{
			ClientID:     conf.ClientID,
			ClientSecret: conf.ClientSecret,
			RedirectURL:  conf.RedirectURL,
			Scopes:       conf.Scopes,
			Endpoint:     githuboauth.Endpoint,
		}
		return oauthConf.AuthCodeURL(state, oauth2.AccessTypeOnline), nil
	case ProviderQQ:
		values := url.Values{}
		values.Set("response_type", "code")
		values.Set("client_id", conf.ClientID)
		values.Set("redirect_uri", conf.RedirectURL)
		values.Set("state", state)
		if len(conf.Scopes) > 0 {
			values.Set("scope", strings.Join(conf.Scopes, ","))
		}
		return "https://graph.qq.com/oauth2.0/authorize?" + values.Encode(), nil
	case ProviderWeChat:
		values := url.Values{}
		values.Set("appid", conf.ClientID)
		values.Set("redirect_uri", conf.RedirectURL)
		values.Set("response_type", "code")
		values.Set("scope", strings.Join(conf.Scopes, ","))
		values.Set("state", state)
		return "https://open.weixin.qq.com/connect/qrconnect?" + values.Encode() + "#wechat_redirect", nil
	default:
		return "", fmt.Errorf("unsupported provider")
	}
}
