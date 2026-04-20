package provider

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/config"
)

// WeChatClient provides WeChat authorize URL generation and keeps finish flow disabled for phase 1.
// WeChatClient 提供微信授权地址生成，并在一期保持 finish 流程关闭。
type WeChatClient struct {
	Conf config.OAuthProviderConf
}

// NewWeChatClient creates a WeChat provider client.
// NewWeChatClient 创建微信 provider 客户端。
func NewWeChatClient(conf config.OAuthProviderConf) *WeChatClient {
	return &WeChatClient{Conf: conf}
}

// Name returns the provider name.
// Name 返回 provider 名称。
func (c *WeChatClient) Name() string {
	return auth.ProviderWeChat
}

// Enabled reports whether the provider is configured and enabled.
// Enabled 判断 provider 是否配置启用。
func (c *WeChatClient) Enabled() bool {
	return c != nil && c.Conf.Enabled
}

// LoginReady reports whether the WeChat login flow is fully implemented.
// LoginReady 判断微信登录流程是否已完整实现。
func (c *WeChatClient) LoginReady() bool {
	return false
}

// RedirectURL returns the configured redirect URL.
// RedirectURL 返回配置中的回调地址。
func (c *WeChatClient) RedirectURL() string {
	return strings.TrimSpace(c.Conf.RedirectURL)
}

// BuildAuthorizeURL builds the WeChat authorize URL.
// BuildAuthorizeURL 构建微信授权地址。
func (c *WeChatClient) BuildAuthorizeURL(state string) (string, error) {
	if c == nil {
		return "", fmt.Errorf("wechat client is nil")
	}

	values := url.Values{}
	values.Set("appid", c.Conf.ClientID)
	values.Set("redirect_uri", c.Conf.RedirectURL)
	values.Set("response_type", "code")
	values.Set("scope", strings.Join(c.Conf.Scopes, ","))
	values.Set("state", state)

	return "https://open.weixin.qq.com/connect/qrconnect?" + values.Encode() + "#wechat_redirect", nil
}
