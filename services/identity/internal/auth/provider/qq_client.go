package provider

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/config"
)

// QQClient provides QQ authorize URL generation and keeps finish flow disabled for phase 1.
// QQClient 提供 QQ 授权地址生成，并在一期保持 finish 流程关闭。
type QQClient struct {
	Conf config.OAuthProviderConf
}

// NewQQClient creates a QQ provider client.
// NewQQClient 创建 QQ provider 客户端。
func NewQQClient(conf config.OAuthProviderConf) *QQClient {
	return &QQClient{Conf: conf}
}

// Name returns the provider name.
// Name 返回 provider 名称。
func (c *QQClient) Name() string {
	return auth.ProviderQQ
}

// Enabled reports whether the provider is configured and enabled.
// Enabled 判断 provider 是否配置启用。
func (c *QQClient) Enabled() bool {
	return c != nil && c.Conf.Enabled
}

// LoginReady reports whether the QQ login flow is fully implemented.
// LoginReady 判断 QQ 登录流程是否已完整实现。
func (c *QQClient) LoginReady() bool {
	return false
}

// RedirectURL returns the configured redirect URL.
// RedirectURL 返回配置中的回调地址。
func (c *QQClient) RedirectURL() string {
	return strings.TrimSpace(c.Conf.RedirectURL)
}

// BuildAuthorizeURL builds the QQ authorize URL.
// BuildAuthorizeURL 构建 QQ 授权地址。
func (c *QQClient) BuildAuthorizeURL(state string) (string, error) {
	if c == nil {
		return "", fmt.Errorf("qq client is nil")
	}

	values := url.Values{}
	values.Set("response_type", "code")
	values.Set("client_id", c.Conf.ClientID)
	values.Set("redirect_uri", c.Conf.RedirectURL)
	values.Set("state", state)
	if len(c.Conf.Scopes) > 0 {
		values.Set("scope", strings.Join(c.Conf.Scopes, ","))
	}

	return "https://graph.qq.com/oauth2.0/authorize?" + values.Encode(), nil
}
