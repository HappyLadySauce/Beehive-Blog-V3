package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/auth"
	"github.com/HappyLadySauce/Beehive-Blog-V3/services/identity/internal/config"
)

const (
	weChatAuthorizeURL = "https://open.weixin.qq.com/connect/qrconnect"
	weChatTokenURL     = "https://api.weixin.qq.com/sns/oauth2/access_token"
	weChatUserInfoURL  = "https://api.weixin.qq.com/sns/userinfo"
)

// WeChatClient provides WeChat website-login OAuth capabilities.
// WeChatClient 提供微信网站登录 OAuth 能力。
type WeChatClient struct {
	HTTPClient   *http.Client
	AuthorizeURL string
	TokenURL     string
	UserInfoURL  string
	Conf         config.OAuthProviderConf
}

// NewWeChatClient creates a WeChat OAuth client with production defaults.
// NewWeChatClient 使用生产默认值创建微信 OAuth 客户端。
func NewWeChatClient(conf config.OAuthProviderConf) *WeChatClient {
	return &WeChatClient{
		HTTPClient:   &http.Client{Timeout: 10 * time.Second},
		AuthorizeURL: weChatAuthorizeURL,
		TokenURL:     weChatTokenURL,
		UserInfoURL:  weChatUserInfoURL,
		Conf:         conf,
	}
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

// LoginReady reports whether the provider has a complete login implementation.
// LoginReady 判断 provider 是否具备完整登录实现。
func (c *WeChatClient) LoginReady() bool {
	return c != nil &&
		c.Enabled() &&
		strings.TrimSpace(c.Conf.ClientID) != "" &&
		strings.TrimSpace(c.Conf.ClientSecret) != "" &&
		strings.TrimSpace(c.Conf.RedirectURL) != ""
}

// RedirectURL returns the configured redirect URL.
// RedirectURL 返回配置中的回调地址。
func (c *WeChatClient) RedirectURL() string {
	return strings.TrimSpace(c.Conf.RedirectURL)
}

// BuildAuthorizeURL builds the WeChat QR connect authorize URL.
// BuildAuthorizeURL 构建微信扫码登录授权地址。
func (c *WeChatClient) BuildAuthorizeURL(state string) (string, error) {
	query := url.Values{}
	query.Set("appid", strings.TrimSpace(c.Conf.ClientID))
	query.Set("redirect_uri", strings.TrimSpace(c.Conf.RedirectURL))
	query.Set("response_type", "code")
	query.Set("scope", "snsapi_login")
	query.Set("state", strings.TrimSpace(state))

	return strings.TrimRight(c.AuthorizeURL, "?") + "?" + query.Encode() + "#wechat_redirect", nil
}

// ExchangeCode exchanges a WeChat authorization code for an access token.
// ExchangeCode 使用微信授权码交换 access token。
//
// The WeChat website-login token endpoint does not require redirect_uri on this step.
// 微信网站登录的 token 接口在此步骤不要求传递 redirect_uri。
func (c *WeChatClient) ExchangeCode(ctx context.Context, code, _ string) (*AccessToken, error) {
	tokenResp, err := c.exchangeTokenResponse(ctx, code)
	if err != nil {
		return nil, err
	}

	return &AccessToken{
		Token:        strings.TrimSpace(tokenResp.AccessToken),
		RefreshToken: strings.TrimSpace(tokenResp.RefreshToken),
		OpenID:       strings.TrimSpace(tokenResp.OpenID),
		UnionID:      strings.TrimSpace(tokenResp.UnionID),
		Scope:        strings.TrimSpace(tokenResp.Scope),
	}, nil
}

// FetchProfile fetches the authenticated WeChat profile and normalizes it.
// FetchProfile 拉取已授权微信用户资料并完成标准化。
func (c *WeChatClient) FetchProfile(ctx context.Context, accessToken *AccessToken) (*Profile, []byte, error) {
	if accessToken == nil || strings.TrimSpace(accessToken.Token) == "" || strings.TrimSpace(accessToken.OpenID) == "" {
		return nil, nil, fmt.Errorf("wechat access token context is incomplete")
	}

	query := url.Values{}
	query.Set("access_token", strings.TrimSpace(accessToken.Token))
	query.Set("openid", strings.TrimSpace(accessToken.OpenID))
	query.Set("lang", "zh_CN")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.UserInfoURL+"?"+query.Encode(), nil)
	if err != nil {
		return nil, nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("wechat user info API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var rawMap map[string]any
	if err := json.Unmarshal(body, &rawMap); err != nil {
		return nil, nil, err
	}
	if errCode, ok := rawMap["errcode"].(float64); ok && int(errCode) != 0 {
		return nil, nil, fmt.Errorf("wechat user info API returned error %.0f", errCode)
	}

	openID, _ := rawMap["openid"].(string)
	unionID, _ := rawMap["unionid"].(string)
	if strings.TrimSpace(openID) == "" {
		openID = strings.TrimSpace(accessToken.OpenID)
	}
	if strings.TrimSpace(unionID) == "" {
		unionID = strings.TrimSpace(accessToken.UnionID)
	}
	nickname, _ := rawMap["nickname"].(string)
	headImgURL, _ := rawMap["headimgurl"].(string)

	subject := strings.TrimSpace(unionID)
	subjectType := "unionid"
	if subject == "" {
		subject = strings.TrimSpace(openID)
		subjectType = "openid"
	}
	if subject == "" || strings.TrimSpace(nickname) == "" {
		return nil, nil, fmt.Errorf("wechat profile is missing required fields")
	}

	raw, err := json.Marshal(rawMap)
	if err != nil {
		return nil, nil, err
	}

	return &Profile{
		Subject:          subject,
		SubjectType:      subjectType,
		DisplayName:      strings.TrimSpace(nickname),
		AvatarURL:        stringPtr(headImgURL),
		UnionID:          stringPtr(unionID),
		OpenID:           stringPtr(openID),
		RawProfile:       raw,
		ProviderClientID: stringPtr(strings.TrimSpace(c.Conf.ClientID)),
		RequestedScopes:  scopeStringPtr(accessToken.Scope),
	}, raw, nil
}

type weChatTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid"`
	ErrCode      int64  `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
}

func (c *WeChatClient) exchangeTokenResponse(ctx context.Context, code string) (*weChatTokenResponse, error) {
	query := url.Values{}
	query.Set("appid", strings.TrimSpace(c.Conf.ClientID))
	query.Set("secret", strings.TrimSpace(c.Conf.ClientSecret))
	query.Set("code", strings.TrimSpace(code))
	query.Set("grant_type", "authorization_code")

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.TokenURL+"?"+query.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("wechat token API returned status %d", resp.StatusCode)
	}

	var payload weChatTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
		return nil, err
	}
	if payload.ErrCode != 0 {
		return nil, fmt.Errorf("wechat token API returned error %d: %s", payload.ErrCode, payload.ErrMsg)
	}
	if strings.TrimSpace(payload.AccessToken) == "" || strings.TrimSpace(payload.OpenID) == "" {
		return nil, fmt.Errorf("wechat token API returned incomplete token payload")
	}

	return &payload, nil
}
